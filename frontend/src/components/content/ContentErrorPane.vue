<script setup>
import IconButton from "@/components/common/IconButton.vue";
import Refresh from "@/components/icons/Refresh.vue";
import ErrorExplain from "@/components/content_value/ErrorExplain.vue";
import {
  capitalize,
  includes,
  isEmpty,
  map,
  split,
  uniqBy,
  isString,
  isNumber,
  isNull,
  isUndefined,
} from "lodash";
import { useThemeVars, NTag } from "naive-ui";
import useBrowserStore from "stores/browser.js";
import useConfigStore from "stores/config.js";
import useConnectionStore from "stores/connections.js";
import usePreferencesStore from "stores/preferences.js";
import { computed, h, nextTick, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";

const themeVars = useThemeVars();

const browserStore = useBrowserStore();
const configStore = useConfigStore();
const connectionStore = useConnectionStore();
const preferencesStore = usePreferencesStore();
const i18n = useI18n();
const data = reactive({
  loading: false,
  server: "",
  options: [],
  selectedOptions: [],
  history: [],
  results: [],
});
const validateFilters = (value) => {
  if (isEmpty(value)) {
    $message.error("filter can't be null");
    return;
  }
  if (!includes(AvailableFilters, capitalize(value))) {
    $message.error("Please select at least one valid filter from []");
    return;
  }
};

const filterServerOption = computed(() => {
  const serverSet = uniqBy(data.history, "server");
  const options = map(serverSet, ({ server }) => ({
    label: server,
    value: server,
  }));
  options.splice(0, 0, {
    label: "common.all",
    value: "",
  });
  return options;
});

const filterOptions = computed(() => {
  // const filters = uniqBy(data.options);
  const options = map(data.options, (item) => ({
    label: item,
    value: item,
  }));
  return options;
});
const tableRef = ref(null);
const handleNamespaceName = (nsname) => {
  if (isEmpty(nsname)) {
    return "N/A";
  }
  if (!nsname.includes("/")) {
    return nsname;
  }
  const [ns, name] = split(nsname, "/");
  return h(
    "div",
    {
      style: {
        marginRight: "6px",
      },
      type: "info",
      bordered: false,
    },
    [
      h(
        NTag,
        {
          style: {
            marginRight: "6px",
          },
          type: "info",
          bordered: false,
        },
        {
          default: () => ns,
        },
      ),
      h(
        "div",
        {
          style: {
            marginRight: "6px",
          },
          type: "info",
          bordered: false,
        },
        {
          default: () => name,
        },
      ),
    ],
  );
};
const columns = computed(() => [
  {
    title: () => i18n.t("error.kind"),
    key: "kind",
    // defaultSortOrder: "ascend",
    // sorter: "default",
    width: 50,
    fixed: "left",
    // align: "center",
    // titleAlign: "center",
    // render: ({ timestamp }, index) => {
    //   return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss");
    // },
    render(row) {
      return h("div", { class: "error-kind" }, row.kind);
    },
  },
  {
    title: () => i18n.t("error.name"),
    key: "name",
    width: 80,
    fixed: "left",
    // align: "center",
    // titleAlign: "center",
    // ellipsis: {
    //   tooltip: true,
    // },
    render(row) {
      return handleNamespaceName(row.name);
    },
  },
  {
    title: () => i18n.t("error.parent"),
    key: "parent",
    width: 50,
    fixed: "left",
    // align: "center",
    // titleAlign: "center",
    // ellipsis: {
    //   tooltip: true,
    // },
    render(row) {
      return handleNamespaceName(row.parentObject);
    },
  },
  {
    title: () => i18n.t("error.error"),
    key: "error",
    width: 100,
    fixed: "left",
    // align: "left",
    // titleAlign: "center",
    render(row) {
      const tags = row.error.map((object) => {
        return h(
          "div",
          // NTag,
          {
            style: {
              marginRight: "6px",
              color: "#CB5064",
            },
            type: "info",
            bordered: false,
          },
          {
            default: () => object.Text,
          },
        );
      });
      return tags;
    },
  },
  {
    title: () => i18n.t("error.advice"),
    key: "details",
    width: 150,
    fixed: "left",
    render(row) {
      console.log("row.details.error: ", row.details.error);
      return h(ErrorExplain, { data: row.details });
      // return row.details || "N/A";
    },
  },
]);
const pagination = {
  pageSize: 10,
};
const analyze = async () => {
  if (isEmpty(connectionStore.currentCluster)) {
    $message.warning(i18n.t("error.no_cluster_selected"));
    return;
  }
  // feat: provide selectable options fetched from backend
  console.log("analyze: ", data.selectedOptions);
  try {
    await nextTick();
    data.loading = true;
    console.log("analyze ai:", preferencesStore.ai);
    const resp = await configStore.analyze(
      connectionStore.currentCluster,
      preferencesStore.ai.backend,
      preferencesStore.ai.model,
      data.selectedOptions,
      preferencesStore.ai.explain,
      preferencesStore.ai.aggregate,
      false,
    );
    if (!resp.success) {
      data.results = [];
      // fix: mistake the response.msg to response.message
      $message.error(resp.msg);
    } else {
      data.results = resp.data;
      if (isEmpty(data.results)) {
        $message.success("No error found");
      }
    }
  } finally {
    data.loading = false;
    await nextTick();
    // tableRef.value?.scrollTo({ position: "bottom" });
  }
  console.log(data.results);
};

const loadResources = async () => {
  try {
    await nextTick();
    data.loading = true;
    const list = await connectionStore.getAvailableResources();
    data.options = list || [];
    console.log("data.options: ", data.options);
  } finally {
    data.loading = false;
    await nextTick();
    // tableRef.value?.scrollTo({ position: "bottom" });
  }
};

const cleanHistory = async () => {
  $dialog.warning(i18n.t("error.confirm_clean_log"), async () => {
    try {
      data.loading = true;
      const success = await browserStore.cleanCmdHistory();
      if (success) {
        data.history = [];
        await nextTick();
        tableRef.value?.scrollTo({ position: "top" });
        $message.success(i18n.t("dialogue.handle_succ"));
      }
    } finally {
      data.loading = false;
    }
  });
};

const onSelectedItemUpdate = (keys) => {
  // fix: be compatiable with other possible types
  if (isString(keys)) {
    keys = [keys];
  } else if (isNumber(keys)) {
    keys = [keys.toString()];
  } else if (isNull(keys)) {
    keys = [];
  } else if (isUndefined(keys)) {
    keys = [];
  }
  data.selectedOptions = keys;
  console.log("selectedUpdate: ", data.selectedOptions);
};

defineExpose({
  refresh: loadResources,
});
</script>
<!-- feat: support multiple choice on select -->
<template>
  <div
    class="content-log content-container content-value fill-height flex-box-v"
  >
    <n-h3>{{ $t("error.title") }}</n-h3>
    <n-form :disabled="data.loading" class="flex-item" inline>
      <n-form-item :label="$t('error.filter_resource')">
        <n-select
          :consistent-menu-width="false"
          :options="filterOptions"
          multiple
          filterable
          @update:value="onSelectedItemUpdate"
          style="min-width: 200px"
        />
      </n-form-item>
      <!-- <n-form-item :label="$t('error.filter_keyword')">
        <n-input v-model:value="data.filters" clearable placeholder="" />
      </n-form-item> -->
      <n-form-item label="&nbsp;">
        <icon-button
          :icon="Refresh"
          border
          t-tooltip="error.refresh"
          @click="analyze"
        />
      </n-form-item>
      <!-- <n-form-item label="&nbsp;">
        <icon-button
          :icon="Delete"
          border
          t-tooltip="error.clean_log"
          @click="cleanHistory"
        />
      </n-form-item> -->
    </n-form>
    <n-data-table
      ref="tableRef"
      :columns="columns"
      :pagination="pagination"
      :data="data.results"
      :loading="data.loading"
      class="flex-item-expand"
      flex-height
      virtual-scroll
      striped
      :scroll-x="1800"
    />
  </div>
</template>

<style lang="scss" scoped>
@import "@/styles/content";
</style>
