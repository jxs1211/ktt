<template>
  <n-tabs
    v-model:value="activeTab"
    type="card"
    :addable="true"
    :closable="closable"
    tab-style="min-width: 80px;"
    @close="closeCli"
    @add="addCli"
    @update:value="connectToCli"
  >
    <n-tab-pane
      v-for="session in sessions"
      :key="`${session.address}:${session.port} - ${session.cmds}`"
      :name="`${session.address}:${session.port} - ${session.cmds}`"
      :tab="`${session.address}:${session.port} - ${session.cmds}`"
      class="flex-item-expand"
    >
      <!-- {{ `${session.address}:${session.port} - ${session.cmds}` }} -->
      <!-- <ContentCli
        v-if="connectingCliSession"
        ref="cliRef"
        class="flex-item-expand"
      /> -->
      <Cli
        :address="session.address"
        :port="session.port"
        :cmds="session.cmds"
      />
    </n-tab-pane>
  </n-tabs>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted, reactive } from "vue";
import { NTabs, NTabPane, NForm, NFormItem, useMessage } from "naive-ui";
import IconButton from "@/components/common/IconButton.vue";
import Refresh from "@/components/icons/Refresh.vue";
import Add from "@/components/icons/Add.vue";
import Cli from "@/components/content_value/Cli.vue";
import { GetSessionsByClusterName } from "wailsjs/go/db/DBService.js";
import useDialogStore from "stores/dialog.js";
import useConnectionStore from "stores/connections.js";
import { useSessionStore } from "@/stores/session.js";
import { getPlatform } from "@/utils/platform.js";

const sessionStore = useSessionStore();
const dialogStore = useDialogStore();
const connectionStore = useConnectionStore();

const randomPort = (start, end, excludes = []) => {
  if (start < 1 || end > 65535 || start > end) {
    throw new Error("Invalid port range");
  }
  let port;
  do {
    port = Math.floor(Math.random() * (end - start + 1)) + start;
  } while (excludes.includes(port));
  return port;
};

const getShell = () => {
  const platform = getPlatform().toLowerCase();
  if (platform === "win32" || platform === "win") {
    return "powershell.exe";
  } else if (platform === "linux") {
    return "bash";
  } else if (platform === "darwin") {
    return "zsh";
  }
  return "bash";
};
const addCli = () => {
  // dialogStore.openNewCliDialog();
  // sessionStore.formValue.cluster_name = connectionStore.currentCluster;
  // sessionStore.formValue.address = "127.0.0.1";
  // sessionStore.formValue.port = randomPort(12110, 22110);
  // sessionStore.formValue.cmds = getShell();
  const newCli = {
    address: "127.0.0.1",
    port: String(randomPort(12110, 22110)),
    cmds: getShell()
  }
  sessions.value.push(newCli)
  activeTab.value = newCli
};


const loading = ref(false);
const connectingCliSession = ref(false);
const cliRef = ref(null);
const initVal = {
  address: "127.0.0.1",
  port: String(randomPort(12110, 22110)),
  cmds: getShell(),
};
const activeTab = ref(initVal);
const sessions = ref([initVal]);
const refresh = async () => {
  loading.value = true;
  const resp = await GetSessionsByClusterName(connectionStore.currentCluster);
  if (!resp.success) {
    console.error("get session by cluster name failed: ", resp.msg);
    message.error("Failed to refresh sessions");
  } else {
    sessions.value = resp.data;
  }
  loading.value = false;
};
const connectToCli = (clusterName) => {
  console.log("param: ", clusterName)
  const session = sessions.value.find((s) => s.cluster_name === clusterName);
  if (!session) return;

  cliRef.value?.doStartTerminal(
    session.cluster_name,
    session.address,
    session.port,
    session.cmds,
  );
  connectingCliSession.value = true;
};


const closable = computed(() => {
  return sessions.value.length > 1;
});
const closeCli = async (name) => {
  console.log("close: ", name);
  return;
  const nameIndex = sessions.value.findIndex((panelName) => panelName === name);
  if (!~nameIndex)
    return;
  panels.splice(nameIndex, 1);
  if (name === valueRef.value) {
    valueRef.value = panels[Math.min(nameIndex, panels.length - 1)];
  }
  const resp = await cliRef.value?.doCloseTerminal(session);
  if (!resp.success) {
    console.error("del cli failed: ", resp.msg);
    message.error("Failed to delete CLI session");
    return;
  }
  message.success("CLI session deleted");
};
onMounted(() => {
  // refresh();
});

onUnmounted(() => {
  // Any cleanup if needed
});
</script>

<style scoped>
.cli-bar {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.cli-bar :deep(.n-tabs) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.cli-bar :deep(.n-tabs-nav) {
  flex-shrink: 0;
}

.cli-bar :deep(.n-tab-pane) {
  flex: 1;
  overflow: auto;
}
</style>
