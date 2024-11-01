<script setup>
import IconButton from "@/components/common/IconButton.vue";
import Refresh from "@/components/icons/Refresh.vue";
import ErrorExplain from "@/components/content_value/ErrorExplain.vue";
import Terminal from "@/components/content/Terminal.vue";
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
  filter,
} from "lodash";
import { useThemeVars, NTag, NButton } from "naive-ui";
import useBrowserStore from "stores/browser.js";
import useConfigStore from "stores/config.js";
import useConnectionStore from "stores/connections.js";
import usePreferencesStore from "stores/preferences.js";
import { useSessionStore } from "stores/session.js";
import { watch, computed, h, nextTick, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import useDialogStore from "../../stores/dialog";
import CliBar from "@/components/content_value/CliBar.vue";
import emitter from '@/utils/eventBus';

const cliBarVisible = ref(false);
const cliBarHeight = ref("50%");
const selectedRow = ref(null);
// const themeVars = useThemeVars();
const dialogStore = useDialogStore();
const configStore = useConfigStore();
const connectionStore = useConnectionStore();
const preferencesStore = usePreferencesStore();
const sessionStore = useSessionStore();
const i18n = useI18n();
const data = reactive({
  loading: false,
  // fix filter data loading using table's loading reative data for displaying loading process
  filterResourceLoading: false,
  filterNamespaceLoading: false,
  server: "",
  options: [],
  namespaces: [],
  selectedOptions: [],
  selectedNSOption: "",
  selectedResourceName: "",
  history: [],
  results: [],
});
// const validateFilters = (value) => {
//   if (isEmpty(value)) {
//     $message.error("filter can't be null");
//     return;
//   }
//   if (!includes(AvailableFilters, capitalize(value))) {
//     $message.error("Please select at least one valid filter from []");
//     return;
//   }
// };

// const filterServerOption = computed(() => {
//   const serverSet = uniqBy(data.history, "server");
//   const options = map(serverSet, ({ server }) => ({
//     label: server,
//     value: server,
//   }));
//   options.splice(0, 0, {
//     label: "common.all",
//     value: "",
//   });
//   return options;
// });
const filterNamespaceOptions = computed(() => {
  const options = map(data.namespaceOptions, (item) => ({
    label: item,
    value: item,
  }));
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
      // return h(ErrorExplain, { data: row.details });
      return h(
        NButton,
        { type: "primary", size: "small", onClick: () => debugWithAI(row) },
        { default: () => "debug with AI" },
      );
      // return row.details || "N/A";
    },
  },
]);
const debugWithAI = (row) => {
  if (preferencesStore.ai.enable && !dialogStore.preferencesDialogVisible) {
    selectedRow.value = row;
    // Generate session name
    const sessionName = `${row.kind}: ${row.name}`;
    
    // Prepare initial prompt from row data
    const initialPrompt = `Analyze and suggest fixes for the following Kubernetes issue:
    Resource Kind: ${row.kind}
    Resource Name: ${row.name}
    Parent Object: ${row.parentObject}
    Errors: ${row.error.map(e => e.Text).join('\n')}`;

    // Show CliBar
    cliBarVisible.value = true;
    sessionStore.debugRowMsg = initialPrompt;

    // Pass data to CliBar component
    nextTick(() => {
      // Trigger chat with initial prompt
      emitter.emit('start-debug-session', {
        sessionName,
        initialPrompt,
        row
      });
      console.log("sending to start debug: ", initialPrompt)

    });
  } else if (!preferencesStore.ai.enable) {
    dialogStore.openPreferencesDialog("ai");
  }
};
// Add drag handling
const isDragging = ref(false);
const startY = ref(0);
const startHeight = ref(0);

const startDragging = (e) => {
  isDragging.value = true;
  startY.value = e.clientY;
  startHeight.value = parseInt(cliBarHeight.value);
  
  // Add event listeners
  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', stopDragging);
};

const onDrag = (e) => {
  if (!isDragging.value) return;
  
  const deltaY = startY.value - e.clientY;
  const containerHeight = document.querySelector('.debug-container').clientHeight;
  const newHeight = Math.min(Math.max(100, startHeight.value + deltaY), containerHeight - 100);
  
  cliBarHeight.value = `${newHeight}px`;
};

const stopDragging = () => {
  isDragging.value = false;
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', stopDragging);
};
// Add style computation for data table
const tableStyle = computed(() => ({
  height: cliBarVisible.value ? `calc(100% - ${cliBarHeight.value})` : '100%',
  transition: 'height 0.3s ease'
}));
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
    data.loading = true;
    await nextTick();
    console.log("analyze ai:", preferencesStore.ai);
    const backend = preferencesStore.getBackend(preferencesStore.ai.backend);
    const resp = await configStore.analyze(
      connectionStore.currentCluster,
      backend.name,
      backend.model,
      backend.baseUrl,
      data.selectedOptions,
      data.selectedNSOption,
      // set to false, by far the design is to debug issue 
      // by fetching AI sepeartely using backend GetCompletion API
      // preferencesStore.ai.explain,
      false,
      preferencesStore.ai.aggregate,
      false,
    );
    if (!resp.success) {
      data.results = [];
      // fix: mistake the response.msg to response.message
      console.log("here");
      $message.error(resp.msg);
    } else {
      data.results = resp.data;
      if (!isEmpty(data.selectedResourceName)) {
        data.results = filter(data.results, (ele) => {
          const [_, name] = split(ele.name, "/");
          return name === data.selectedResourceName;
        });
      }
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

const refreshFiltersOptions = async () => {
  try {
    data.filterResourceLoading = true;
    data.filterNamespaceLoading = true;
    await nextTick();
    const resp = await connectionStore.getAvailableResources();
    if (!resp.success) {
      console.warn("get available filtered resources failed: ", resp.msg);
      return;
    }
    data.options = resp.data;
    console.log("data.options: ", data.options);
    const resp2 = await connectionStore.getNamespaces();
    if (!resp2.success) {
      console.warn("get namespaces failed: ", resp2.msg);
    }
    data.namespaceOptions = resp2.data;
    console.log("data.namespaceOptions: ", data.namespaceOptions);
  } finally {
    data.filterResourceLoading = false;
    data.filterNamespaceLoading = false;
    await nextTick();
    // tableRef.value?.scrollTo({ position: "bottom" });
  }
};
const clearSelectedOptions = () => {
  data.selectedOptions = []; // Clear selected resource types
  data.selectedNSOption = ""; // Clear selected namespace
  data.selectedResourceName = "";
  data.results = [];
};
// const cleanHistory = async () => {
//   $dialog.warning(i18n.t("error.confirm_clean_log"), async () => {
//     try {
//       data.loading = true;
//       const success = await browserStore.cleanCmdHistory();
//       if (success) {
//         data.history = [];
//         await nextTick();
//         tableRef.value?.scrollTo({ position: "top" });
//         $message.success(i18n.t("dialogue.handle_succ"));
//       }
//     } finally {
//       data.loading = false;
//     }
//   });
// };

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
const onSelectedNSItemUpdate = (key) => {
  console.log("---->", typeof key, key);
  data.selectedNSOption = key;
};

defineExpose({
  refreshFiltersOptions,
  clearSelectedOptions,
});

// feat: support to reset error table after cluster changed
watch(
  () => connectionStore.currentCluster,
  async (newVal, oldVal) => {
    data.results = [];
    // data.selectedNSOption = "";
    // data.selectedOptions = [];
  },
);
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
          v-model:value="data.selectedOptions"
          :consistent-menu-width="false"
          :options="filterOptions"
          multiple
          filterable
          clearable
          :loading="data.filterResourceLoading"
          @update:value="onSelectedItemUpdate"
          style="min-width: 200px"
        />
      </n-form-item>
      <n-form-item :label="$t('error.filter_namespace')">
        <n-select
          v-model:value="data.selectedNSOption"
          :consistent-menu-width="false"
          :options="filterNamespaceOptions"
          filterable
          clearable
          :loading="data.filterNamespaceLoading"
          @update:value="onSelectedNSItemUpdate"
          style="min-width: 200px"
        />
      </n-form-item>
      <n-form-item :label="$t('error.filter_resource_name')">
        <n-input
          v-model:value="data.selectedResourceName"
          type="text"
          placeholder="resource name"
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
    <div class="table-container" :style="tableStyle">
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
        <!-- Add CliBar component -->
    <!-- Draggable splitter -->
    <div 
      v-if="cliBarVisible"
      class="splitter"
      @mousedown="startDragging"
    ></div>

    <div 
      v-if="cliBarVisible"
      class="cli-bar-container"
      :style="{ height: cliBarHeight }"
    >
      <cli-bar
        v-model:visible="cliBarVisible"
        v-model:height="cliBarHeight"
      />
    </div>
    />
  </div>
</template>

<style lang="scss" scoped>
@import "@/styles/content";

.content-log {
  position: relative;
  overflow: hidden;
}

.debug-container {
  position: relative;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.table-container {
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
  min-height: 100px;
}

.splitter {
  height: 6px;
  background: var(--n-border-color);
  cursor: row-resize;
  position: relative;
  
  &:hover {
    background: var(--n-primary-color);
  }
  
  &::before {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 30px;
    height: 2px;
    background: currentColor;
    border-radius: 1px;
  }
}

.cli-bar-container {
  position: relative;
  min-height: 100px;
  transition: height 0.3s ease;
  
  :deep(.cli-bar) {
    height: 100%;
  }
}

:deep(.n-data-table) {
  height: 100%;
  flex: 1;
}

// Disable text selection while dragging
body.dragging {
  user-select: none;
}
</style>