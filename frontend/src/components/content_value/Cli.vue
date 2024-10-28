<script setup>
import { isEmpty } from "lodash";
import { onMounted, onUnmounted, ref } from "vue";
import {
  CloseTerminal2,
  StartTerminal2,
} from "wailsjs/go/cli/TerminalService.js";
import { EventsOff, EventsOn } from "wailsjs/runtime/runtime.js";

import useConnectionStore from "stores/connections.js";
import useTabStore from "stores/tab.js";

const props = defineProps({
  name: String,
  address: {
    type: String,
    // default: "127.0.0.1",
    required: true,
  },
  port: {
    type: String,
    required: true,
  },
  cmds: {
    type: String,
    required: true,
  },
  activated: Boolean,
});

const terminalUrl = ref("");
const iframeRef = ref(null);
const connectionStore = useConnectionStore();
const tabStore = useTabStore();
const randomPort = () => {
  return Math.floor(Math.random() * (13110 - 12110 + 1)) + 12110;
}; // Generates a random port number between 12110 and 22110
const serverOptions = {
  address: "127.0.0.1",
  port: String(randomPort()),
  cmds: "zsh",
};

const isValidUrl = (url) => {
  const pattern =
    /^(https?:\/\/((([0-9]{1,3}\.){3}[0-9]{1,3})|([a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+))(:\d+)?(\/[^\s]*)?)$/; // Updated regex for validating URL
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

const doStartTerminal = async (address, port, cmds) => {
  console.log("doStartTerminal2: ", address, port, cmds);
  EventsOn("terminal2:url", (url) => {
    console.log("url: ", url);
    terminalUrl.value = url;
  });
  try {
    const resp = await StartTerminal2(address, port, cmds);
    if (!resp.success) {
      console.error("failed to start terminal", resp.msg);
      return;
    }
    console.log(resp);
  } catch (error) {
    console.error("Failed to start terminal:", error);
  }
};
const doCloseTerminal = async (session) => {
  EventsOff("terminal2:url");
  terminalUrl.value = null;
  console.log(
    "start to tear down: ",
    session.address,
    session.port,
    session.cmds,
  );
  return await CloseTerminal2(session.address, session.port, session.cmds);
};
const refreshTerminal = async (address, port) => {};
const getTerminal = async (address, port) => {};
const getAllTerminals = async () => {};
onMounted(() => {
  console.log(
    "cli bar's cli onMounted: ",
    props.address,
    props.port,
    props.cmds,
  );
  doStartTerminal(props.address, props.port, props.cmds);
});
onUnmounted(() => {
  console.log("console cli tab onUnmounted");
});
defineExpose({
  doStartTerminal,
  doCloseTerminal,
});
</script>

<template>
  <div class="terminal-container">
    <n-spin v-if="isEmpty(terminalUrl)" class="spinner">
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
