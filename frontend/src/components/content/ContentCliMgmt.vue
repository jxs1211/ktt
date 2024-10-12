<script setup>
import IconButton from "@/components/common/IconButton.vue";
import Refresh from "@/components/icons/Refresh.vue";
import Add from "@/components/icons/Add.vue";
import ContentCli from "@/components/content_value/ContentCli.vue";
import Left from "@/components/icons/Left.vue"
import EditCliDialog from "@/components/dialogs/EditCliDialog.vue";
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
import useDialogStore from "stores/dialog.js";
import useConfigStore from "stores/config.js";
import useConnectionStore from "stores/connections.js";
import usePreferencesStore from "stores/preferences.js";
import { useSessionStore } from '@/stores/session.js';
import { getPlatform } from "@/utils/platform.js";
import { watch, computed, h, nextTick, reactive, ref, onMounted, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import {
  // CreateSession,
  // DeleteSession,
  GetSessionsByClusterName,
} from 'wailsjs/go/db/DBService.js';

// const themeVars = useThemeVars();
const sessionStore = useSessionStore();
const dialogStore = useDialogStore();
const connectionStore = useConnectionStore();
// const configStore = useConfigStore();
// const ferencesStore = usePreferencesStore();
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
  // results: [{"cluster_name": "wave-loadtest-others-us-east-1.us-east-1.eksctl.io", "address":"0.0.0.0", "port": "1211", "cmds": "zsh"}],
});
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
    title: () => i18n.t("cli_mgmt.cluster_name"),
    key: "ClusterName",
    // // defaultSortOrder: "ascend",
    // // sorter: "default",
    // width: 20,
    // fixed: "left",
    // // align: "center",
    // // titleAlign: "center",
    // // render: ({ timestamp }, index) => {
    // //   return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss");
    // // },
    render(row) {
      return h(
        "div", 
        {style: {marginRight: "6px"}, type: "info", bordered: false},
        row.cluster_name
      );
    },
  },
  {
    title: () => i18n.t("cli_mgmt.address"),
    key: "Address",
    // width: 10,
    // fixed: "left",
    render(row) {
      return h(
        "div",
        {style: {marginRight: "6px"}, type: "info", bordered: false},
        {default: () => row.address},
      );
    },
  },
  {
    title: () => i18n.t("cli_mgmt.port"),
    key: "Port",
    // width: 10,
    // fixed: "left",
    render(row) {
      return h(
        "div",
        {style: {marginRight: "6px"}, type: "info", bordered: false},
        {default: () => row.port},
      );
    },
  },
  {
    title: () => i18n.t("cli_mgmt.cmds"),
    key: "Cmds",
    // width: 100,
    // fixed: "left",
    render(row) {
      return h(
        "div",
        {style: {marginRight: "6px"}, type: "info", bordered: false},
        {default: () => row.cmds},
      );
    },
  },
  {
    title: () => i18n.t("cli_mgmt.actions"),
    key: 'Actions',
    // width: 30,
    // fixed: "right",
    render(row) {
      return h('div', { style: { display: 'flex', gap: '10px' } }, [
        h(NButton, { type: 'primary', size: 'small', onClick: () => connectToCli(row) }, 'Connect'),
        h(NButton, { type: 'info', size: 'small', onClick: () => delCli(row) }, 'Delete'),
        h(NButton, { type: 'info', size: 'small', onClick: () => editCli(row) }, 'Edit'),
      ]);
    },
  },
]);
const connectingCliSession = ref(false);
const selectedSession = reactive({
  clusterName: "",
  address: "",
  port: "",
  cmds: "",
}); // To hold the session data
const connectToCli = (session) => {
  // Logic to open ContentCli.vue with the session details
  // This could involve setting a shared state or using a router
  console.log('Connecting to CLI:', session);
  cliRef.value?.doStartTerminal(
    session.cluster_name,
    session.address,
    session.port,
    session.cmds
  )
  // console.log("selectedSession: ", selectedSession.clusterName, selectedSession.address, selectedSession.port, selectedSession.cmds)
  connectingCliSession.value = true; // 
  // Example: Set the session in a store or emit an event
};

// const deleteCli = async (session) => {
//   const resp = await DeleteSession(session.id);
//   if (!resp.success) {
//     console.error("Failed to delete session:", resp.msg);
//     $message.error("failed to delete session: ", resp.msg);
//     return;
//   }
//   // Remove the session from the store
//   sessionStore.results = sessionStore.results.filter(s => s.id !== session.id);
//   console.log("Session deleted:", session);
//   $message.success("session deleted: ", session)
// };
const pagination = {
  pageSize: 10,
};
// Generates a random port number between 12110 and 22110
const startPort = 12110;
const endPort = 22110;
const cliRef = ref(null);
const editCliDialogRef = ref(null);
const randomPort = (start, end, excludes = []) => {
	if (start < 1 || end > 65535 || start > end) {
		throw new Error("Invalid port range");
	}
	let port;
	do {
		port = Math.floor(Math.random() * (end - start + 1)) + start;
	} while (excludes.includes(port)); // Keep generating until found a port not in excludes
	return port;
};
// kttodo: generate cmds by resource type
const getShell = () => {
	const platform = getPlatform().toLowerCase();
	console.log("platform: ", platform);
	if (platform === 'win32' || platform === "win") {
		return 'powershell.exe'; // or 'cmd.exe' for PowerShell
	} else if (platform === 'linux') {
		return 'bash'; // Default shell for most Linux distributions
	} else if (platform === 'darwin') {
		return 'zsh'; // Default shell for macOS (as of Catalina)
	}
	return 'bash'; // Default to bash if the platform is unknown
};
const backToList = async () => {
  connectingCliSession.value = false
}
const addCli = async () => {
  dialogStore.openNewCliDialog();
  sessionStore.formValue.cluster_name = connectionStore.currentCluster;
  sessionStore.formValue.address = "127.0.0.1" 
  sessionStore.formValue.port = randomPort(startPort, endPort) 
  sessionStore.formValue.cmds = getShell()
  console.log("addCli")
};
const delCli = async (row) => {
  const resp = await cliRef.value?.doCloseTerminal(row)
  if (!resp.success) {
    console.error("del cli failed: ", resp.msg)
    $message.error("del cli failed: " + resp.msg)
    return
  }
  console.log("do close terminal resp: ", resp)
  const resp2 = await GetSessionsByClusterName(row.cluster_name);
  if (!resp2.success) {
    console.error("refresh sessions by cluster name failed: ", resp2.msg)
    $message.error("refresh sessions by cluster name failed: " + resp2.msg)
    return
  }
  console.log("get data after del cli: ", resp2.data)
  sessionStore.setResults(resp2.data)
  connectingCliSession.value = false
  $message.success("del cli ok")
};
const editCli = (row) => {
  // open edit window
  // edit form value
  // save and refresh data
  console.log("edit cli row: ", row)
  dialogStore.openEditCliDialog();
  // sessionStore.formValue.id = row.id
  // sessionStore.formValue.cluster_name = row.cluster_name;
  // sessionStore.formValue.address = row.address;
  // sessionStore.formValue.port = row.port;
  // sessionStore.formValue.cmds = row.cmds;
  editCliDialogRef.value?.doUpdateTerminal(row)

  // const resp = await doEditTerminal(row)
}
const refreshFiltersOptions = async () => {
  try {
    data.filterResourceLoading = true;
    data.filterNamespaceLoading = true;
    await nextTick();
    const resp = await connectionStore.getAvailableResources();
    if (!resp.success) {
      console.warn("get available filtered resources failed: ", resp.msg)
      return
    }
    data.options = resp.data;
    console.log("data.options: ", data.options);
    const resp2 = await connectionStore.getNamespaces()
    if (!resp2.success) {
      console.warn("get namespaces failed: ", resp2.msg)
    }
    data.namespaceOptions = resp2.data;
    console.log("data.namespaceOptions: ", data.namespaceOptions);
  } finally {
    data.filterResourceLoading = false
    data.filterNamespaceLoading = false
    await nextTick();
    // tableRef.value?.scrollTo({ position: "bottom" });
  }
};
const clearSelectedOptions = () => {
  data.selectedOptions = []; // Clear selected resource types
  data.selectedNSOption = ""; // Clear selected namespace
  data.selectedResourceName = "";
  sessionStore.results = [];
};
// const onSelectedItemUpdate = (keys) => {
//   // fix: be compatiable with other possible types
//   if (isString(keys)) {
//     keys = [keys];
//   } else if (isNumber(keys)) {
//     keys = [keys.toString()];
//   } else if (isNull(keys)) {
//     keys = [];
//   } else if (isUndefined(keys)) {
//     keys = [];
//   }
//   data.selectedOptions = keys;
//   console.log("selectedUpdate: ", data.selectedOptions);
// };
// const onSelectedNSItemUpdate = (key) => {
//   console.log("---->", typeof(key), key)
//   data.selectedNSOption = key;
// };

// defineExpose({
//   refreshFiltersOptions,
//   clearSelectedOptions,
// });

// feat: support to reset error table after cluster changed
watch(
  () => connectionStore.currentCluster,
  async (newVal, oldVal) => {
    sessionStore.results = [];
    // data.selectedNSOption = "";
    // data.selectedOptions = [];
  },
);
onMounted(() => {
  console.log("content cli mgmt on mounted")
  refresh(connectionStore.currentCluster)
});
onUnmounted(() => {
  console.log("content cli mgmt on unmounted")
}
);
const refresh = async (cluster) => {
  const resp = await GetSessionsByClusterName(cluster)
  if (!resp.success) {
    console.log("get session by cluster name failed: ", resp.msg)
    return
  }
  sessionStore.setResults(resp.data)
  console.log("resp data is ", resp.data)
};
</script>
<template>
  <div
    class="content-log content-container content-value fill-height flex-box-v"
  >
    <n-form :disabled="data.loading" class="flex-item" inline>
      <!-- <n-form-item :label="$t('error.filter_resource_name')" v-show="!connectingCliSession">
        <n-input
          v-model:value="data.selectedResourceName" 
          type="text" 
          placeholder="resource name" />
      </n-form-item> -->
      <n-form-item label="&nbsp;" v-show="!connectingCliSession">
        <icon-button
          :icon="Refresh"
          border
          t-tooltip="error.refresh"
          @click="() => refresh(connectionStore.currentCluster)"
        />
      </n-form-item>
      <n-form-item label="&nbsp;" v-show="!connectingCliSession">
        <icon-button
          :icon="Add"
          border
          @click="addCli"
        />
      </n-form-item>
      <n-form-item label="&nbsp;" v-show="connectingCliSession">
        <icon-button
          :icon="Left"
          border
          @click="backToList"
        />
      </n-form-item>
    </n-form>
    <n-data-table
      v-show="!connectingCliSession"
      ref="tableRef"
      :columns="columns"
      :pagination="pagination"
      :data="sessionStore.results"
      :loading="data.loading"
      class="flex-item-expand"
      flex-height
      virtual-scroll
      striped
      :scroll-x="180"
    />
    <ContentCli
      v-show="connectingCliSession"
      ref="cliRef"
    />
    <EditCliDialog v-show="false" ref="editCliDialogRef" />
  </div>
</template>

<style lang="scss" scoped>
@import "@/styles/content";
</style>
