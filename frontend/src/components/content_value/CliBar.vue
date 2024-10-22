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
    class="cli-bar"
  >
    <n-tab-pane
      v-for="(session, index) in sessions"
      :key="sessionKey(session)"
      :name="sessionKey(session)"
      :tab="sessionKey(session)"
      class="flex-item-expand"
    >
      <Terminal
        :ref="
          (el) => {
            if (el) cliRefs[index] = el;
          }
        "
        :enableChat="true"
        :initPrompt="'init-prompt'"
        :address="session.address"
        :port="session.port"
        :cmds="session.cmds"
      />
    </n-tab-pane>
  </n-tabs>
</template>

<script setup>
import Terminal from "@/components/content/Terminal.vue";
import { extraTheme } from "@/utils/extra_theme.js";
import { computed, ref, onMounted, onUnmounted, reactive } from "vue";
import { NTabs, NTabPane, NForm, NFormItem, useMessage } from "naive-ui";
import IconButton from "@/components/common/IconButton.vue";
import Refresh from "@/components/icons/Refresh.vue";
import Add from "@/components/icons/Add.vue";
import ResizeableWrapper from "@/components/common/ResizeableWrapper.vue";
import Cli from "@/components/content_value/Cli.vue";
import { GetSessionsByClusterName } from "wailsjs/go/db/DBService.js";
import useDialogStore from "stores/dialog.js";
import usePreferencesStore from "stores/preferences.js";
import useConnectionStore from "stores/connections.js";
import { useSessionStore } from "@/stores/session.js";
import { getPlatform } from "@/utils/platform.js";

const chatWidth = ref(300); // Initial width of the chat component
const sessionStore = useSessionStore();
const dialogStore = useDialogStore();
const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const sessionKey = (session) =>
  `${session.address}:${session.port} - ${session.cmds}`;
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
const exThemeVars = computed(() => {
  return extraTheme(prefStore.isDark);
});
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
  const newCli = {
    address: "127.0.0.1",
    port: String(randomPort(12110, 22110)),
    cmds: getShell(),
  };
  sessions.value.push(newCli);
  activeTab.value = sessionKey(newCli);
};

const loading = ref(false);
const connectingCliSession = ref(false);
const cliRef = ref(null);
const cliRefs = ref([]);
const activeTab = ref();
const sessions = ref();
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
const connectToCli = (key) => {
  console.log("update tab param: ", key);
  return;
  const session = sessions.value.find((s) => sessionKey(s) === key);
  if (!session) {
    $message.error("Session not found: ", key);
    return;
  }

  cliRef.value?.doStartTerminal(session.address, session.port, session.cmds);
  // connectingCliSession.value = true;
};

const closable = computed(() => {
  return true;
  // return sessions.value.length > 1;
});
const closeCli = async (key) => {
  console.log("close: ", key);
  const sessIndex = sessions.value.findIndex((s) => sessionKey(s) === key);
  if (!~sessIndex) {
    console.error("Session not found: ", key);
    return;
  }
  // const resp = await cliRef.value?.doCloseTerminal(sessions.value[sessIndex]);
  const resp = await cliRefs.value[sessIndex]?.doCloseTerminal(
    sessions.value[sessIndex],
  );
  if (!resp.success) {
    console.error("Failed to delete CLI session: ", resp.msg);
    $message.error("Failed to delete CLI session: ", resp.msg);
    return;
  }

  sessions.value.splice(sessIndex, 1);
  if (key === activeTab.value) {
    activeTab.value =
      sessions.value[Math.min(sessIndex, sessions.value.length - 1)];
  }
  $message.success("CLI session deleted");
};
onMounted(() => {
  // Initialize with a default session
  const initialSession = {
    address: "127.0.0.1",
    port: String(randomPort(12110, 22110)),
    cmds: getShell(),
  };
  sessions.value = [initialSession];
  activeTab.value = sessionKey(initialSession);
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

.cli-bar :deep(.n-tabs-content) {
  flex: 1;
  overflow: hidden;
}

.cli-bar :deep(.n-tab-pane) {
  height: 100%;
}
</style>
