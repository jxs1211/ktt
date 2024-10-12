<script setup>
import { isEmpty } from 'lodash';
import { onMounted, onUnmounted, ref } from 'vue';
import { CloseTerminal, StartTerminal } from 'wailsjs/go/cli/TerminalService.js';
import { EventsOff, EventsOn } from 'wailsjs/runtime/runtime.js';

import useConnectionStore from "stores/connections.js";
import useTabStore from "stores/tab.js";

const props = defineProps({
  name: String,
  activated: Boolean,
});

const terminalUrl = ref('');
const iframeRef = ref(null);
const connectionStore = useConnectionStore();
const tabStore = useTabStore();
const randomPort = () => { return Math.floor(Math.random() * (13110 - 12110 + 1)) + 12110; }; // Generates a random port number between 12110 and 22110
const serverOptions = {
  address: "127.0.0.1",
  port: String(randomPort()),
  cmds: "zsh",
};

const isValidUrl = (url) => {
  const pattern = /^(https?:\/\/((([0-9]{1,3}\.){3}[0-9]{1,3})|([a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+))(:\d+)?(\/[^\s]*)?)$/; // Updated regex for validating URL
  return pattern.test(url); // Test the URL against the pattern
};
const split = (url) => {
  if (!isValidUrl(url)) {
    console.error("Invalid URL format: ", url); // Log error for invalid URL
    return null; // Return null or handle the error as needed
  }

  const urlParts = new URL(url); // Create a URL object
  const address = urlParts.hostname; // Get the hostname (address)
  console.log("parsed addr: ", address);
  const port = urlParts.port; // Get the port
  return { address, port }; // Return as an object
};

const doStartTerminal = async (cluster, address, port, cmds) => {
  console.log("doStartTerminal: ", cluster, address, port, cmds)
  EventsOn('terminal:url', (url) => {
    console.log("url: ", url)
    terminalUrl.value = url;
  });
  try {
    const resp = await StartTerminal(cluster, address, port, cmds);
    if (!resp.success) {
      console.error('failed to start terminal', resp.msg)
      return
    }
    console.log(resp)
  } catch (error) {
    console.error('Failed to start terminal:', error);
  }
}
const doCloseTerminal = async (row) => {
  EventsOff('terminal:url');
  terminalUrl.value = null;
  console.log("start to tear down: ", terminalUrl.value, row.id, row.cluster_name, row.address, row.port, row.cmds);
  return await CloseTerminal(
    row.id, row.cluster_name, row.address, row.port, row.cmds);
  // if (!success) {
  //   console.error("close terminal failed: ", msg)
  //   return
  // }
  // console.log("finished closing terminal");
};
// const refreshTerminal = async (address, port) => {
// };
// const getTerminal = async (address, port) => {
// };
// const getAllTerminals = async () => {
// };
onMounted(() => {
  console.log("console cli tab onMounted: ", connectionStore.currentCluster)
  // doStartTerminal(props.clusterName, props.address, props.port, props.cmds);
});
onUnmounted(() => {
  console.log("console cli tab onUnmounted")
});
// watch(
//   () => connectionStore.currentCluster,
//   async (newVal, oldVal) => {
//     data.results = [];
//     // data.selectedNSOption = "";
//     // data.selectedOptions = [];
//   },
// );
// watch(
//   () => connectionStore.currentCluster,
//   (newVal, oldVal) => {
//     if (isEmpty(oldVal)){
//       console.log("first selected a cluster")
//       return
//     }
//     if (isEmpty(terminalUrl)) {
//       console.log("terminalUrl is empty")
//       return
//     }
//     const { address, port } = split(terminalUrl.value);
//     if (isEmpty(address) || isEmpty(port)) {
//       console.log("address or port is empty: ", address, port)
//       return
//     }
//     doCloseTerminal(connectionStore.currentCluster, address, port, "zsh");
//     console.log("tear down for current cluster changed")
//   }
// );
// watch(
//   () => tabStore.nav,
//   (newVal, oldVal) => {
//     if (oldVal === "browser" && newVal != "browser") {
//       const { address, port } = split(terminalUrl.value);
//       doCloseTerminal(connectionStore.currentCluster, address, port, "zsh");
//       console.log("tear down for nav changed")    
//     }
//   }
// );
// watch(
//   () => tabStore.currentSubTab,
//   (newVal, oldVal) => {
//     if (oldVal === "cli" && newVal != "cli") {
//       const { address, port } = split(terminalUrl.value);
//       doCloseTerminal(address, port);
//     } else if (oldVal != "cli" && newVal === "cli") {
//       doStartTerminal(serverOptions.address, serverOptions.port, serverOptions.cmds);
//     }
//   }
// );
defineExpose({
  doStartTerminal,
  doCloseTerminal,
})
</script>

<template>
  <div class="terminal-container">
    <n-spin v-if="isEmpty(terminalUrl)" class="spinner" >
      <!-- Add your custom spinner here -->
    </n-spin>
    <iframe
      v-else
      :src="terminalUrl"
      ref="iframeRef"
      class="terminal-iframe"
      frameborder="0"
    ></iframe>
  </div>
</template>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
.spinner {
  height: 100%;
  width: 100%;
  position: relative;
}
.terminal-iframe {
  width: 100%;
  height: 100%;
  border: none;
}
</style>