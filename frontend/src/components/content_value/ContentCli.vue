<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { EventsOn, EventsOff } from 'wailsjs/runtime/runtime.js';
import { StartTerminal, CloseTerminal } from 'wailsjs/go/cli/TerminalService.js';

import useTabStore from "stores/tab.js";
import useConnectionStore from "stores/connections.js";

const props = defineProps({
  name: String,
  activated: Boolean,
});

const terminalUrl = ref('');
const iframeRef = ref(null);
const connectionStore = useConnectionStore();
const tabStore = useTabStore();
onMounted(async () => {
  try {
    await StartTerminal();
  } catch (error) {
    console.error('Failed to start terminal:', error);
  }

  EventsOn('terminal:url', (url) => {
    console.log("url: ", url)
    terminalUrl.value = url;
  });
});
const tearDown = () => {
  EventsOff('terminal:url');
  terminalUrl.value = null;
  CloseTerminal();
};
onUnmounted(() => {
  console.log("terminal onUnmounted")
});
// watch(
//   () => connectionStore.currentCluster,
//   async (newVal, oldVal) => {
//     data.results = [];
//     // data.selectedNSOption = "";
//     // data.selectedOptions = [];
//   },
// );
watch(
  () => tabStore.currentSubTab,
  (newVal, oldVal) => {
    if (oldVal == "cli" && newVal != "cli") {
      tearDown();
      console.log("tear down")    
    }
  }
);
</script>

<template>
  <div class="terminal-container">
    <n-spin v-if="terminalUrl == ''" class="spinner" >
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