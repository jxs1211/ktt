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
import useDialog, { ConnDialogType } from "stores/dialog";
import Close from "@/components/icons/Close.vue";
import useConfigStore from "stores/config.js";
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
const configStore = useConfigStore();
const browserStore = useBrowserStore();
const connectionStore = useConnectionStore();
const i18n = useI18n();

const editName = ref("");
const content = ref("");
const onLoadConfig = async () => {
  const result = await configStore.loadConfig(content.value);
  if (result.success) {
    validateResult.value = "";
    connectionStore.clusters = result.data.map((_, index) => {
      return {
        label: result.data[index],
        key: `${index}`,
      };
    });
    console.log(connectionStore.clusters);
  } else {
    validateResult.value = result.msg || $t("dialogue.connection.test_fail");
  }
  if (result.success) {
    $message.success(i18n.t("dialogue.handle_succ"));
    onClose();
  }
};
const onTestConnection = async () => {
  const clusters = await configStore.getClusters();
  console.log(clusters);
  validating.value = true;
  const result = await configStore.testConnection(content.value);
  if (result.success) {
    validateResult.value = "";
  } else {
    validateResult.value = result || $t("dialogue.connection.test_fail");
  }
  validating.value = false;
};
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

const resetData = () => {
  validating.value = false;
  validateResult.value = null;
};
const onClose = () => {
  dialogStore.closeConfigDialog();
};
watch(
  () => dialogStore.configDialogVisible,
  (visible) => {
    if (visible) {
      resetData();
    }
  },
);
</script>

<template>
  <n-modal
    v-model:show="dialogStore.configDialogVisible"
    :closable="false"
    :mask-closable="false"
    :on-after-leave="resetData"
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
          v-model:value="content"
          type="textarea"
          placeholder="config content here"
        />
      </n-space>
      <!-- validation result alert-->
      <n-alert
        v-if="showValidateResult"
        :on-close="() => (validateResult = '')"
        :title="
          isEmpty(validateResult) ? '' : $t('dialogue.connection.test_fail')
        "
        :type="isEmpty(validateResult) ? 'success' : 'error'"
        closable
      >
        <template v-if="isEmpty(validateResult)">{{
          $t("dialogue.connection.test_succ")
        }}</template>
        <template v-else>{{ validateResult }}</template>
      </n-alert>
      <!-- <template #action> -->

      <!-- </template> -->
    </n-spin>

    <template #action>
      <div class="flex-item-expand">
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          :loading="validating"
          @click="onTestConnection"
        >
          {{ $t("dialogue.connection.test") }}
        </n-button>
      </div>
      <div class="flex-item n-dialog__action">
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          @click="onClose"
        >
          {{ $t("common.cancel") }}
        </n-button>
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          type="primary"
          @click="onLoadConfig"
        >
          {{ $t("common.confirm") }}
          <!-- {{ isEditMode ? $t("preferences.general.update") : $t("common.confirm") }} -->
        </n-button>
      </div>
    </template>
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
