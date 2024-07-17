<script setup>
import {
  every,
  get,
  includes,
  isEmpty,
  map,
  reject,
  sortBy,
  toNumber,
  trim,
} from "lodash";
import { computed, nextTick, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import {
  ListSentinelMasters,
  TestConnection,
} from "wailsjs/go/services/connectionService.js";
import useDialog, { ConnDialogType } from "stores/dialog";
import Close from "@/components/icons/Close.vue";
import useConnectionStore from "stores/connections.js";
import FileOpenInput from "@/components/common/FileOpenInput.vue";
import { KeyViewType } from "@/consts/key_view_type.js";
import { useThemeVars } from "naive-ui";
import useBrowserStore from "stores/browser.js";
import Delete from "@/components/icons/Delete.vue";
import Add from "@/components/icons/Add.vue";
import IconButton from "@/components/common/IconButton.vue";

/**
 * Dialog for new or edit connection
 */

const themeVars = useThemeVars();
const dialogStore = useDialog();
const connectionStore = useConnectionStore();
const browserStore = useBrowserStore();
const i18n = useI18n();

const editName = ref("");
const generalForm = ref(null);
const generalFormRules = () => {
  const requiredMsg = i18n.t("dialogue.field_required");
  const illegalChars = ["/", "\\"];
  return {
    name: [
      { required: true, message: requiredMsg, trigger: "input" },
      {
        validator: (rule, value) => {
          return every(illegalChars, (c) => !includes(value, c));
        },
        message: i18n.t("dialogue.illegal_characters"),
        trigger: "input",
      },
    ],
    addr: { required: true, message: requiredMsg, trigger: "input" },
    defaultFilter: { required: true, message: requiredMsg, trigger: "input" },
    keySeparator: { required: true, message: requiredMsg, trigger: "input" },
  };
};
const isEditMode = computed(() => dialogStore.connType === ConnDialogType.EDIT);
const closingConnection = computed(() => {
  if (isEmpty(editName.value)) {
    return false;
  }
  return browserStore.isConnected(editName.value);
});

const groupOptions = computed(() => {
  const options = map(connectionStore.groups, (group) => ({
    label: group,
    value: group,
  }));
  options.splice(0, 0, {
    label: "dialogue.connection.no_group",
    value: "",
  });
  return options;
});

const dbFilterList = ref([]);
const onUpdateDBFilterType = (t) => {
  if (t !== "none") {
    // set default filter index if empty
    if (isEmpty(dbFilterList.value)) {
      dbFilterList.value = ["0"];
    }
  }
};

const aliasPair = ref([
  /*{ db: 0, alias: '' }*/
]);
const onCreateAlias = () => {
  return {
    db: 0,
    alias: "",
  };
};
const onUpdateAlias = () => {
  const val = reject(aliasPair.value, (v) => v == null || isEmpty(v.alias));
  const result = {};
  for (const elem of val) {
    result[elem.db] = elem.alias;
  }
  generalForm.value.alias = result;
};

watch(
  () => dbFilterList.value,
  (list) => {
    const dbList = map(list, (item) => {
      const idx = toNumber(item);
      return isNaN(idx) ? 0 : idx;
    });
    generalForm.value.dbFilterList = sortBy(dbList);
  },
  { deep: true },
);

const sshLoginType = computed(() => {
  return get(generalForm.value, "ssh.loginType", "pwd");
});

const loadingSentinelMaster = ref(false);
const masterNameOptions = ref([]);
const onLoadSentinelMasters = async () => {
  try {
    loadingSentinelMaster.value = true;
    const { success, data, msg } = await ListSentinelMasters(generalForm.value);
    if (!success || isEmpty(data)) {
      $message.error(msg || "list sentinel master fail");
    } else {
      const options = [];
      for (const m of data) {
        options.push({
          label: m["name"],
          value: m["name"],
        });
      }

      // select default names
      if (!isEmpty(options)) {
        generalForm.value.sentinel.master = options[0].value;
      }
      masterNameOptions.value = options;
    }
  } catch (e) {
    $message.error(e.message);
  } finally {
    loadingSentinelMaster.value = false;
  }
};

const tab = ref("general");
const validating = ref(false);
const validateResult = ref(null);
const showValidateResult = computed(() => {
  return !validating.value && validateResult.value != null;
});
const predefineColors = ref([
  "",
  "#F75B52",
  "#F7A234",
  "#F7CE33",
  "#4ECF60",
  "#348CF7",
  "#B270D3",
]);
const generalFormRef = ref(null);
const advanceFormRef = ref(null);

const onSaveConnection = async () => {
  // validate general form
  await generalFormRef.value?.validate((err) => {
    if (err) {
      nextTick(() => (tab.value = "general"));
    }
  });

  // validate advance form
  await advanceFormRef.value?.validate((err) => {
    if (err) {
      nextTick(() => (tab.value = "advanced"));
    }
  });

  // trim addr by network type
  if (get(generalForm.value, "network", "tcp") === "unix") {
    generalForm.value.network = "unix";
    generalForm.value.addr = "";
    generalForm.value.port = 0;
    generalForm.value.sock = trim(generalForm.value.sock);
  } else {
    generalForm.value.network = "tcp";
    generalForm.value.sock = "";
    generalForm.value.addr = trim(generalForm.value.addr);
  }

  // trim advance data
  if (get(generalForm.value, "dbFilterType", "none") === "none") {
    generalForm.value.dbFilterList = [];
  }

  // trim ssl data
  if (!!!generalForm.value.ssl.enable) {
    generalForm.value.ssl = {};
  }

  // trim ssh login data
  if (!!generalForm.value.ssh.enable) {
    switch (generalForm.value.ssh.loginType) {
      case "pkfile":
        generalForm.value.ssh.password = "";
        break;
      default:
        generalForm.value.ssh.pkFile = "";
        generalForm.value.ssh.passphrase = "";
        break;
    }
  } else {
    // ssh disabled, reset to default value
    generalForm.value.ssh = {};
  }

  // trim sentinel data
  if (!!!generalForm.value.sentinel.enable) {
    generalForm.value.sentinel = {};
  }

  // trim cluster data
  if (!!!generalForm.value.cluster.enable) {
    generalForm.value.cluster = {};
  }

  // trim proxy data
  if (generalForm.value.proxy.type !== 2) {
    generalForm.value.proxy.schema = "";
    generalForm.value.proxy.addr = "";
    generalForm.value.proxy.port = 0;
    generalForm.value.proxy.auth = false;
    generalForm.value.proxy.username = "";
    generalForm.value.proxy.password = "";
  } else if (!generalForm.value.proxy.auth) {
    generalForm.value.proxy.username = "";
    generalForm.value.proxy.password = "";
  }

  // store new connection
  const { success, msg } = await connectionStore.saveConnection(
    isEditMode.value ? editName.value : null,
    generalForm.value,
  );
  if (!success) {
    $message.error(msg);
    return;
  }

  $message.success(i18n.t("dialogue.handle_succ"));
  onClose();
};

const resetForm = () => {
  console.log("resetForm called");
  // generalForm.value = connectionStore.newDefaultConnection()
  // generalFormRef.value?.restoreValidation()
  // testing.value = false
  // testResult.value = null
  // tab.value = 'general'
  // loadingSentinelMaster.value = false
};

watch(
  () => dialogStore.connDialogVisible,
  (visible) => {
    if (visible) {
      resetForm();
      editName.value = get(dialogStore.connParam, "name", "");
      generalForm.value =
        dialogStore.connParam || connectionStore.newDefaultConnection();
      dbFilterList.value = map(
        generalForm.value.dbFilterList,
        (item) => item + "",
      );
      generalForm.value.ssh.loginType =
        generalForm.value.ssh.loginType || "pwd";
      // update alias display
      const alias = get(generalForm.value, "alias", {});
      const pairs = [];
      for (const db in alias) {
        pairs.push({ db: parseInt(db), alias: alias[db] });
      }
      aliasPair.value = pairs;
      generalForm.value.proxy.auth = !isEmpty(generalForm.value.proxy.username);
    }
  },
);

const onTestConnection = async () => {
  testResult.value = "";
  testing.value = true;
  let result = "";
  try {
    const { success = false, msg } = await TestConnection(generalForm.value);
    if (!success) {
      result = msg;
    }
  } catch (e) {
    result = e.message;
  } finally {
    testing.value = false;
  }

  if (!isEmpty(result)) {
    testResult.value = result;
  } else {
    testResult.value = "";
  }
};

const onClose = () => {
  dialogStore.closeConfigDialog();
};

const pasteFromClipboard = async () => {
  // url example:
  // rediss://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2
  let opt = {};
  try {
    opt = await connectionStore.parseUrlFromClipboard();
  } catch (e) {
    $message.error(
      i18n.t("dialogue.connection.parse_fail", { reason: e.message }),
    );
    return;
  }
  generalForm.value.network = opt.network || "tcp";
  generalForm.value.name = generalForm.value.addr = opt.addr;
  generalForm.value.port = opt.port;
  generalForm.value.username = opt.username;
  generalForm.value.password = opt.password;
  if (opt.connTimeout > 0) {
    generalForm.value.connTimeout = opt.connTimeout;
  }
  if (opt.execTimeout > 0) {
    generalForm.value.execTimeout = opt.execTimeout;
  }
  const { sslServerName = null } = opt;
  if (sslServerName != null) {
    generalForm.value.ssl.enable = true;
    if (!isEmpty(sslServerName)) {
      generalForm.value.ssl.sni = sslServerName;
    }
  }
  $message.success(i18n.t("dialogue.connection.parse_pass", { url: opt.url }));
};
</script>

<template>
  <n-modal
    v-model:show="dialogStore.configDialogVisible"
    :closable="false"
    :mask-closable="false"
    :on-after-leave="resetForm"
    :show-icon="false"
    :title="
      isEditMode
        ? $t('dialogue.connection.edit_title')
        : $t('dialogue.connection.new_title')
    "
    close-on-esc
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
    @esc="onClose"
  >
    <n-spin :show="closingConnection">
      <n-text>Add your kubeconfig</n-text>
      <n-space vertical>
        <n-input
          v-model:value="value"
          type="textarea"
          placeholder="config content here"
        />
      </n-space>
      <!-- validation result alert-->
      <n-alert
          v-if="showValidateResult"
          :on-close="() => (validateResult = '')"
          :title="isEmpty(validateResult) ? '' : $t('dialogue.connection.test_fail')"
          :type="isEmpty(validateResult) ? 'success' : 'error'"
          closable>
          <template v-if="isEmpty(validateResult)">{{ $t('dialogue.connection.test_succ') }}</template>
          <template v-else>{{ validateResult }}</template>
      </n-alert>
      <!-- <template #action> -->
        <n-button :focusable="false" type="primary" @click="onSaveConnection">
            <!-- {{ isEditMode ? $t('preferences.general.update') : $t('common.confirm') }} -->
            {{ $t('common.confirm') }}
        </n-button>
      <!-- </template> -->
    </n-spin>

    <!-- <template #action>
            <div class="flex-item-expand">
                <n-button :disabled="closingConnection" :focusable="false" :loading="testing" @click="onTestConnection">
                    {{ $t('dialogue.connection.test') }}
                </n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="closingConnection" :focusable="false" @click="pasteFromClipboard">
                    {{ $t('dialogue.connection.parse_url_clipboard') }}
                </n-button>
                <n-button :disabled="closingConnection" :focusable="false" @click="onClose">
                    {{ $t('common.cancel') }}
                </n-button>
                <n-button :disabled="closingConnection" :focusable="false" type="primary" @click="onSaveConnection">
                    {{ isEditMode ? $t('preferences.general.update') : $t('common.confirm') }}
                </n-button>
            </div>
        </template> -->
  </n-modal>
</template>

<style lang="scss" scoped>
.color-preset-item {
  width: 24px;
  height: 24px;
  margin-right: 2px;
  border-width: 3px;
  border-style: solid;
  cursor: pointer;
  border-radius: 50%;
}
</style>
