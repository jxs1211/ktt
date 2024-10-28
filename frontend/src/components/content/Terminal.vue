<template>
  <div class="content-area">
    <resizeable-wrapper
      v-if="enableChat"
      v-model:size="chatWidth"
      :min-size="250"
      :max-size="800"
      class="flex-item"
    >
      <Chat :session="sessionString" class="chat-component" />
    </resizeable-wrapper>
    <Cli
      ref="cliRef"
      :address="address"
      :port="port"
      :cmds="cmds"
      class="cli-component"
    />
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import Cli from "@/components/content_value/Cli.vue";
import Chat from "@/components/content_value/Chat.vue";
import ResizeableWrapper from "@/components/common/ResizeableWrapper.vue";

const props = defineProps({
  enableChat: {
    type: Boolean,
    default: false,
  },
  initPrompt: {
    type: String,
    default: "",
  },
  address: {
    type: String,
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
});

const chatWidth = ref(300);
const cliRef = ref(null);
const sessionString = computed(() => {
  return `${props.address}:${props.port}:${props.cmds}`;
});
const doCloseTerminal = async () => {
  return await cliRef.value?.doCloseTerminal({
    address: props.address,
    port: props.port,
    cmds: props.cmds,
  });
};

defineExpose({
  doCloseTerminal,
});
</script>

<style lang="css" scoped>
.content-area {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.chat-component {
  height: 100%;
  overflow-y: auto;
}

.cli-component {
  flex: 1;
  height: 100%;
  overflow: hidden;
}
</style>
