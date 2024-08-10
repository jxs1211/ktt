<script setup>
import IconButton from "@/components/common/IconButton.vue";
import AddLink from "@/components/icons/AddLink.vue";
import Checked from "@/components/icons/Checked.vue";
import Delete from "@/components/icons/Delete.vue";
import Edit from "@/components/icons/Edit.vue";
import Help from "@/components/icons/Help.vue";
import { typesIconStyle } from "@/consts/support_redis_type.js";
import { joinCommand } from "@/utils/decoder_cmd.js";
import { find, map, sortBy, get } from "lodash";
import { NButton, NEllipsis, NIcon, NSpace, NTooltip } from "naive-ui";
import useDialog from "stores/dialog.js";
import usePreferencesStore from "stores/preferences.js";
import { computed, h, ref, watchEffect, watch } from "vue";
import { useI18n } from "vue-i18n";
import { BrowserOpenURL } from "wailsjs/runtime/runtime.js";

const prefStore = usePreferencesStore();
const modelPlaceholder = ref("Please select a model");
const prevPreferences = ref({});

const localaiModelOptions = computed(() => [
  {
    value: "llama2",
    label: "llama2",
  },
  {
    value: "llama3",
    label: "llama3",
  },
]);
const openaiModelOptions = computed(() => [
  {
    value: "gpt-3.5-turbo",
    label: "gpt-3.5-turbo",
  },
  {
    value: "gpt-4",
    label: "gpt-4",
  },
]);
const tab = ref("general");
const aiProviderTab = ref("localai");
const dialogStore = useDialog();
const i18n = useI18n();
const loading = ref(false);

const initPreferences = async () => {
  try {
    loading.value = true;
    tab.value = dialogStore.preferencesTag || "general";
    await prefStore.loadPreferences();
    prevPreferences.value = {
      general: prefStore.general,
      ai: prefStore.ai,
      editor: prefStore.editor,
      cli: prefStore.cli,
      decoder: prefStore.decoder,
    };
    aiProviderTab.value = prefStore.ai.backend;
  } finally {
    loading.value = false;
  }
};
watchEffect(() => {
  if (dialogStore.preferencesDialogVisible) {
    initPreferences();
  }
});
watch(
  () => aiProviderTab.value,
  (newTabName) => {
    prefStore.onSelectedTab(newTabName);
  },
);
watch(
  () => prefStore.ai.enable,
  (newVal) => {
    if (newVal) {
      if (prefStore.ai.backend === "noopai") {
        // handle the first time enabling ai
        aiProviderTab.value = "localai";
      } else {
        aiProviderTab.value = prefStore.ai.backend;
      }
    }
    prefStore.onSwitch(newVal);
  },
);
const keyOptions = computed(() => {
  const opts = map(typesIconStyle, (v) => ({
    value: v,
    label: "preferences.general.key_icon_style" + v,
  }));
  return sortBy(opts, (o) => o.value);
});

const decoderList = computed(() => {
  const decoder = prefStore.decoder || [];
  const list = [];
  for (const d of decoder) {
    // decode command
    list.push({
      name: d.name,
      auto: d.auto,
      decodeCmd: joinCommand(d.decodePath, d.decodeArgs),
      encodeCmd: joinCommand(d.encodePath, d.encodeArgs),
    });
  }
  return list;
});

const decoderColumns = computed(() => {
  return [
    {
      key: "name",
      title: () => i18n.t("preferences.decoder.decoder_name"),
      width: 120,
      align: "center",
      titleAlign: "center",
    },
    {
      key: "cmd",
      title: () => i18n.t("preferences.decoder.cmd_preview"),
      titleAlign: "center",
      render: ({ decodeCmd, encodeCmd }, index) => {
        return h(
          NSpace,
          {
            vertical: true,
            wrapItem: false,
            wrap: false,
            justify: "center",
            size: 15,
          },
          () => [
            h(
              NEllipsis,
              {},
              {
                default: () => decodeCmd,
                tooltip: () => decodeCmd + "\n\n" + encodeCmd,
              },
            ),
            h(
              NEllipsis,
              {},
              {
                default: () => encodeCmd,
                tooltip: () => decodeCmd + "\n\n" + encodeCmd,
              },
            ),
          ],
        );
      },
    },
    {
      key: "status",
      title: () => i18n.t("preferences.decoder.status"),
      width: 80,
      align: "center",
      titleAlign: "center",
      render: ({ auto }, index) => {
        if (auto) {
          return h(
            NTooltip,
            { delay: 0, showArrow: false },
            {
              default: () => i18n.t("preferences.decoder.auto_enabled"),
              trigger: () => h(NIcon, { component: Checked, size: 16 }),
            },
          );
        }
        return "-";
      },
    },
    {
      key: "action",
      title: () => i18n.t("interface.action"),
      width: 80,
      align: "center",
      titleAlign: "center",
      render: ({ name, auto }, index) => {
        return h(
          NSpace,
          { wrapItem: false, wrap: false, justify: "center", size: "small" },
          () => [
            h(IconButton, {
              icon: Delete,
              tTooltip: "interface.delete_row",
              onClick: () => {
                prefStore.removeCustomDecoder(name);
              },
            }),
            h(IconButton, {
              icon: Edit,
              tTooltip: "interface.edit_row",
              onClick: () => {
                const decoders = prefStore.decoder || [];
                const decoder = find(decoders, { name });
                const { auto, decodePath, decodeArgs, encodePath, encodeArgs } =
                  decoder;
                dialogStore.openDecoderDialog({
                  name,
                  auto,
                  decodePath,
                  decodeArgs,
                  encodePath,
                  encodeArgs,
                });
              },
            }),
          ],
        );
      },
    },
  ];
});

const onOpenPrivacy = () => {
  let helpUrl = "";
  switch (prefStore.currentLanguage) {
    case "zh":
      helpUrl = "https://redis.tinycraft.cc/zh/guide/privacy.html";
      break;
    default:
      helpUrl = "https://redis.tinycraft.cc/guide/privacy.html";
      break;
  }
  BrowserOpenURL(helpUrl);
};

const openDecodeHelp = () => {
  let helpUrl = "";
  switch (prefStore.currentLanguage) {
    case "zh":
      helpUrl = "https://redis.tinycraft.cc/zh/guide/custom-decoder.html";
      break;
    default:
      helpUrl = "https://redis.tinycraft.cc/guide/custom-decoder.html";
      break;
  }
  BrowserOpenURL(helpUrl);
};

const onSavePreferences = async () => {
  const success = await prefStore.savePreferences();
  if (success) {
    if (prefStore.ai.enable) {
      const message = i18n.t("dialogue.handle_succ");
      const backend = prefStore.ai.backend;
      $message.success(`${message}: ${backend} is selected`);
    }
    dialogStore.closePreferencesDialog();
  }
};

const onClose = () => {
  // restore to old preferences
  prefStore.resetToLastPreferences();
  dialogStore.closePreferencesDialog();
};
const localaiModel = computed({
  get: () => prefStore.getBackend("localai")?.model || "",
  set: (value) => {
    const backend = prefStore.getBackend("localai");
    if (backend) {
      backend.model = value;
    }
  },
});
const localaiBaseUrl = computed({
  get: () => prefStore.getBackend("localai")?.baseUrl || "",
  set: (value) => {
    const backend = prefStore.getBackend("localai");
    if (backend) {
      backend.baseUrl = value;
    }
  },
});
const openaiModel = computed({
  get: () => prefStore.getBackend("openai")?.model || "",
  set: (value) => {
    const backend = prefStore.getBackend("openai");
    if (backend) {
      backend.model = value;
    }
  },
});
const openaiApiKey = computed({
  get: () => prefStore.getBackend("openai")?.apiKey || "",
  set: (value) => {
    const backend = prefStore.getBackend("openai");
    if (backend) {
      backend.apiKey = value;
    }
  },
});
</script>

<template>
  <n-modal
    v-model:show="dialogStore.preferencesDialogVisible"
    :auto-focus="false"
    :closable="false"
    :mask-closable="false"
    :show-icon="false"
    :title="$t('preferences.name')"
    close-on-esc
    preset="dialog"
    style="width: 640px"
    transform-origin="center"
    @esc="onClose"
  >
    <!-- FIXME: set loading will slow down appear animation of dialog in linux -->
    <!-- <n-spin :show="loading"> -->
    <n-tabs
      v-model:value="tab"
      animated
      pane-style="min-height: 300px"
      placement="left"
      tab-style="justify-content: right; font-weight: 420;"
      type="line"
    >
      <!-- general pane -->
      <n-tab-pane
        :tab="$t('preferences.general.name')"
        display-directive="show"
        name="general"
      >
        <n-form
          :disabled="loading"
          :model="prefStore.general"
          :show-require-mark="false"
          label-placement="top"
        >
          <n-grid :x-gap="10">
            <n-form-item-gi
              :label="$t('preferences.general.theme')"
              :span="24"
              required
            >
              <n-radio-group
                v-model:value="prefStore.general.theme"
                name="theme"
                size="medium"
              >
                <n-radio-button
                  v-for="opt in prefStore.themeOption"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ $t(opt.label) }}
                </n-radio-button>
              </n-radio-group>
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.language')"
              :span="24"
              required
            >
              <n-select
                v-model:value="prefStore.general.language"
                :options="prefStore.langOption"
                :render-label="
                  ({ label, value }) => (value === 'auto' ? $t(label) : label)
                "
                filterable
              />
            </n-form-item-gi>
            <n-form-item-gi :span="24" required>
              <template #label>
                {{ $t("preferences.general.font") }}
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-icon :component="Help" />
                  </template>
                  <div class="text-block">
                    {{ $t("preferences.font_tip") }}
                  </div>
                </n-tooltip>
              </template>
              <n-select
                v-model:value="prefStore.general.fontFamily"
                :options="prefStore.fontOption"
                :placeholder="$t('preferences.general.font_tip')"
                :render-label="
                  ({ label, value }) => (value === '' ? $t(label) : label)
                "
                filterable
                multiple
                tag
              />
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.font_size')"
              :span="24"
            >
              <n-input-number
                v-model:value="prefStore.general.fontSize"
                :max="65535"
                :min="1"
              />
            </n-form-item-gi>
            <n-form-item-gi :span="12">
              <template #label>
                {{ $t("preferences.general.scan_size") }}
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-icon :component="Help" />
                  </template>
                  <div class="text-block">
                    {{ $t("preferences.general.scan_size_tip") }}
                  </div>
                </n-tooltip>
              </template>
              <n-input-number
                v-model:value="prefStore.general.scanSize"
                :min="1"
                :show-button="false"
                style="width: 100%"
              />
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.key_icon_style')"
              :span="12"
            >
              <n-select
                v-model:value="prefStore.general.keyIconStyle"
                :options="keyOptions"
                :render-label="({ label }) => $t(label)"
              />
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.update')"
              :span="24"
            >
              <n-checkbox v-model:checked="prefStore.general.checkUpdate">
                {{ $t("preferences.general.auto_check_update") }}
              </n-checkbox>
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.privacy')"
              :span="24"
            >
              <n-checkbox v-model:checked="prefStore.general.allowTrack">
                {{ $t("preferences.general.allow_track") }}
                <n-button
                  style="text-decoration: underline"
                  text
                  type="primary"
                  @click="onOpenPrivacy"
                >
                  {{ $t("preferences.general.privacy") }}
                </n-button>
              </n-checkbox>
            </n-form-item-gi>
          </n-grid>
        </n-form>
      </n-tab-pane>
      <!-- ai pane -->
      <n-tab-pane
        :tab="$t('preferences.ai.name')"
        display-directive="show"
        name="ai"
      >
        <n-grid :cols="1">
          <n-form-item-gi label="Enable">
            <n-switch v-model:value="prefStore.ai.enable" />
          </n-form-item-gi>
          <n-form-item-gi v-if="prefStore.ai.enable" label="Provider list">
            <n-tabs
              v-model:value="aiProviderTab"
              animated
              pane-style="min-height: 300px"
              placement="top"
              tab-style="justify-content: right; font-weight: 420;"
              type="line"
            >
              <n-tab-pane
                :tab="$t('preferences.ai.providers.localai.name')"
                display-directive="show"
                name="localai"
              >
                <n-grid :cols="1" :style="{ width: '90%' }">
                  <n-form-item-gi
                    :label="$t('preferences.ai.providers.localai.model')"
                  >
                    <n-select
                      v-model:value="localaiModel"
                      filterable
                      :placeholder="modelPlaceholder"
                      :options="localaiModelOptions"
                    />
                  </n-form-item-gi>
                  <n-form-item-gi
                    :label="$t('preferences.ai.providers.localai.base_url')"
                  >
                    <n-input v-model:value="localaiBaseUrl" type="text" />
                  </n-form-item-gi>
                </n-grid>
              </n-tab-pane>
              <n-tab-pane
                :tab="$t('preferences.ai.providers.openai.name')"
                display-directive="show"
                name="openai"
              >
                <n-grid :cols="1" :style="{ width: '90%' }">
                  <n-form-item-gi
                    :label="$t('preferences.ai.providers.openai.model')"
                  >
                    <n-select
                      v-model:value="openaiModel"
                      filterable
                      :placeholder="modelPlaceholder"
                      :options="openaiModelOptions"
                    />
                  </n-form-item-gi>
                  <n-form-item-gi :label="$t('preferences.ai.api_key')">
                    <n-input v-model:value="openaiApiKey" type="text" />
                  </n-form-item-gi>
                </n-grid>
              </n-tab-pane>
              <!-- <n-tab-pane
                :tab="$t('preferences.ai.providers.azure.name')"
                display-directive="show"
                name="azure"
              >
                <n-grid :cols="1">
                  <n-form-item-gi
                    :label="$t('preferences.ai.providers.azure.model')"
                  >
                    <n-input v-model:value="prefStore.ai.model" type="text" />
                  </n-form-item-gi>
                  <n-form-item-gi :label="$t('preferences.ai.api_key')">
                    <n-input v-model:value="prefStore.ai.apiKey" type="text" />
                  </n-form-item-gi>
                </n-grid>
              </n-tab-pane> -->
            </n-tabs>
          </n-form-item-gi>
        </n-grid>
      </n-tab-pane>
      <!-- editor pane -->
      <n-tab-pane
        :tab="$t('preferences.editor.name')"
        display-directive="show"
        name="editor"
      >
        <n-form
          :disabled="loading"
          :model="prefStore.editor"
          :show-require-mark="false"
          label-placement="top"
        >
          <n-grid :x-gap="10">
            <n-form-item-gi :span="24" required>
              <template #label>
                {{ $t("preferences.general.font") }}
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-icon :component="Help" />
                  </template>
                  <div class="text-block">
                    {{ $t("preferences.font_tip") }}
                  </div>
                </n-tooltip>
              </template>
              <n-select
                v-model:value="prefStore.editor.fontFamily"
                :options="prefStore.fontOption"
                :placeholder="$t('preferences.general.font_tip')"
                :render-label="({ label, value }) => value || $t(label)"
                filterable
                multiple
                tag
              />
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.font_size')"
              :span="24"
            >
              <n-input-number
                v-model:value="prefStore.editor.fontSize"
                :max="65535"
                :min="1"
              />
            </n-form-item-gi>
            <n-form-item-gi
              :show-feedback="false"
              :show-label="false"
              :span="24"
            >
              <n-checkbox v-model:checked="prefStore.editor.showLineNum">
                {{ $t("preferences.editor.show_linenum") }}
              </n-checkbox>
            </n-form-item-gi>
            <n-form-item-gi
              :show-feedback="false"
              :show-label="false"
              :span="24"
            >
              <n-checkbox v-model:checked="prefStore.editor.showFolding">
                {{ $t("preferences.editor.show_folding") }}
              </n-checkbox>
            </n-form-item-gi>
            <n-form-item-gi
              :show-feedback="false"
              :show-label="false"
              :span="24"
            >
              <n-checkbox v-model:checked="prefStore.editor.dropText">
                {{ $t("preferences.editor.drop_text") }}
              </n-checkbox>
            </n-form-item-gi>
            <n-form-item-gi
              :show-feedback="false"
              :show-label="false"
              :span="24"
            >
              <n-checkbox v-model:checked="prefStore.editor.links">
                {{ $t("preferences.editor.links") }}
              </n-checkbox>
            </n-form-item-gi>
          </n-grid>
        </n-form>
      </n-tab-pane>

      <!-- cli pane -->
      <n-tab-pane
        :tab="$t('preferences.cli.name')"
        display-directive="show"
        name="cli"
      >
        <n-form
          :disabled="loading"
          :model="prefStore.cli"
          :show-require-mark="false"
          label-placement="top"
        >
          <n-grid :x-gap="10">
            <n-form-item-gi :span="24" required>
              <template #label>
                {{ $t("preferences.general.font") }}
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-icon :component="Help" />
                  </template>
                  <div class="text-block">
                    {{ $t("preferences.font_tip") }}
                  </div>
                </n-tooltip>
              </template>
              <n-select
                v-model:value="prefStore.cli.fontFamily"
                :options="prefStore.fontOption"
                :placeholder="$t('preferences.general.font_tip')"
                :render-label="({ label, value }) => value || $t(label)"
                filterable
                multiple
                tag
              />
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.general.font_size')"
              :span="24"
            >
              <n-input-number
                v-model:value="prefStore.cli.fontSize"
                :max="65535"
                :min="1"
              />
            </n-form-item-gi>
            <n-form-item-gi
              :label="$t('preferences.cli.cursor_style')"
              :span="24"
            >
              <n-radio-group
                v-model:value="prefStore.cli.cursorStyle"
                name="theme"
                size="medium"
              >
                <n-radio-button
                  v-for="opt in prefStore.cliCursorStyleOption"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ $t(opt.label) }}
                </n-radio-button>
              </n-radio-group>
            </n-form-item-gi>
          </n-grid>
        </n-form>
      </n-tab-pane>

      <!-- custom decoder pane -->
      <n-tab-pane
        :tab="$t('preferences.decoder.name')"
        display-directive="show:lazy"
        name="decoder"
      >
        <n-space vertical>
          <n-space justify="space-between">
            <n-button @click="dialogStore.openDecoderDialog()">
              <template #icon>
                <n-icon :component="AddLink" size="18" />
              </template>
              {{ $t("preferences.decoder.new") }}
            </n-button>
            <n-button @click="openDecodeHelp">
              <template #icon>
                <n-icon :component="Help" size="18" />
              </template>
              {{ $t("preferences.decoder.help") }}
            </n-button>
          </n-space>
          <n-data-table
            :columns="decoderColumns"
            :data="decoderList"
            :single-line="false"
            max-height="350px"
          />
        </n-space>
      </n-tab-pane>
    </n-tabs>
    <!-- </n-spin> -->

    <template #action>
      <div class="flex-item-expand">
        <n-button :disabled="loading" @click="prefStore.restorePreferences">
          {{ $t("preferences.restore_defaults") }}
        </n-button>
      </div>
      <div class="flex-item n-dialog__action">
        <n-button :disabled="loading" @click="onClose">{{
          $t("common.cancel")
        }}</n-button>
        <n-button :disabled="loading" type="primary" @click="onSavePreferences">
          {{ $t("common.save") }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<style lang="scss" scoped>
.inline-form-item {
  padding-right: 10px;
}
</style>
