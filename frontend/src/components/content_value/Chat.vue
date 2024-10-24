<template>
  <div class="chat-container flex-box-v">
    <div class="chat-header">
      <div class="chat-title">CHAT</div>
      <div class="header-controls">
        <n-button quaternary circle size="small" @click="startNewChat">
          <template #icon
            ><n-icon><Add /></n-icon
          ></template>
        </n-button>
        <n-button quaternary circle size="small" @click="showHistory">
          <template #icon
            ><n-icon><Time /></n-icon
          ></template>
        </n-button>
        <n-button quaternary circle size="small" @click="closeChat">
          <template #icon
            ><n-icon><Close /></n-icon
          ></template>
        </n-button>
      </div>
    </div>
    <div class="message-display flex-item">
      <div v-for="message in messages" :key="message.id" class="message">
        <div v-if="message.type === 'text'" class="text-message">
          {{ message.content }}
        </div>
        <code-block
          v-else-if="message.type === 'code'"
          :code="message.content"
          :language="message.language"
        />
      </div>
    </div>
    <div class="chat-input-wrapper">
      <div class="chat-input" ref="chatInputRef">
        <div class="input-wrapper">
          <n-input
            ref="inputRef"
            v-model:value="inputMessage"
            type="textarea"
            placeholder="Add context"
            @keydown.enter.prevent="handleEnterKey"
            @input="adjustTextareaHeight"
          />
        </div>
      </div>
      <div class="command-panel">
        <!-- <div class="left-controls"> -->
        <n-select
          v-model:value="selectedModel"
          :options="modelOptions"
          size="tiny"
          filterable
          clearable
          :render-label="renderModelLabel"
          class="model-select"
        />
        <!-- </div> -->
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, h, nextTick, onMounted } from "vue";
import { NButton, NIcon, NInput, NSelect } from "naive-ui";
import { CaretUp, Add, Time, Close, ChevronDown } from "@vicons/ionicons5";
import CodeBlock from "./CodeBlock.vue";
import { useThemeVars } from "naive-ui";

/**
 * Resizeable component wrapper
 */
const themeVars = useThemeVars();
const messages = ref([]);
const inputMessage = ref("");
// Model selection
const selectedModel = ref("cursor-small");
const modelOptions = [
  {
    label: "cursor-small",
    value: "cursor-small",
    icon: "cursor-icon", // You might want to replace this with actual icon
  },
  {
    label: "cursor-medium",
    value: "cursor-medium",
    icon: "cursor-icon",
  },
];

// Mention selection
const selectedMention = ref(null);
const mentionOptions = [
  {
    label: "@Mention",
    value: "mention",
  },
  // Add more mention options as needed
];

// Modified render label to include caret prefix
const renderModelLabel = (option) => {
  return h(
    "div",
    {
      style: {
        display: "flex",
        alignItems: "center",
        gap: "4px",
      },
    },
    [
      // h(NIcon, { size: 12 }, { default: () => h(ChevronDown) }),
      h("span", { class: "model-label" }, option.label),
    ],
  );
};
// Modified message handling
const chatInputRef = ref(null);
const inputRef = ref(null);

const adjustTextareaHeight = () => {
  if (inputRef.value && chatInputRef.value) {
    const textarea = inputRef.value.$el.querySelector("textarea");
    const chatInput = chatInputRef.value;

    // Reset heights
    textarea.style.height = "auto";

    // Calculate new heights
    const newTextareaHeight = Math.min(
      Math.max(textarea.scrollHeight, 60),
      400,
    ); // Increased max height to 400px
    const newChatInputHeight = newTextareaHeight + 40; // 40px for command panel

    // Apply new heights
    textarea.style.height = `${newTextareaHeight}px`;
    chatInput.style.height = `${newChatInputHeight}px`;
  }
};

const handleEnterKey = (e) => {
  if (e.shiftKey) {
    e.preventDefault();
    const textarea = e.target;
    const cursorPos = textarea.selectionStart;
    const textBeforeCursor = inputMessage.value.substring(0, cursorPos);
    const textAfterCursor = inputMessage.value.substring(cursorPos);

    inputMessage.value = textBeforeCursor + "\n" + textAfterCursor;

    nextTick(() => {
      textarea.selectionStart = cursorPos + 1;
      textarea.selectionEnd = cursorPos + 1;
      adjustTextareaHeight();
    });
  } else {
    sendMessage();
  }
};

// Modify the sendMessage function to reset the textarea height
const sendMessage = () => {
  console.log("Sending message:", inputMessage.value);
  if (inputMessage.value.trim()) {
    messages.value.push({
      id: Date.now(),
      type: "text",
      content: inputMessage.value.trim(),
    });
    inputMessage.value = "";
    nextTick(() => {
      adjustTextareaHeight();
    });
  }
};
const startNewChat = () => {
  messages.value = [];
};

const showHistory = () => {
  // Implement history logic
};

const closeChat = () => {
  // Implement close logic
};
onMounted(() => {
  adjustTextareaHeight();
});
</script>

<style scoped>
.chat-container {
  height: 100%;
  width: 100%;
  background-color: #0000;
  border-left: 1px solid var(--n-border-color);
  box-sizing: border-box;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  position: relative; /* Added this line */
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--n-border-color);
  height: 40px;
  width: 100%;
  box-sizing: border-box;
}

.message-display {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  width: 100%;
  box-sizing: border-box;
  margin-bottom: 96px; /* Added this line */
}

.chat-input-wrapper {
  position: absolute; /* Changed from relative to absolute */
  bottom: 0; /* Added this line */
  left: 0; /* Added this line */
  right: 0; /* Added this line */
  min-height: 96px;
  max-height: 440px;
  border-top: 1px solid var(--n-border-color);
  display: flex;
  flex-direction: column;
  background-color: #0000; /* Added this line */
  z-index: 2; /* Added this line */
}

.chat-input {
  padding: 12px 12px 40px 12px; /* Modified bottom padding */
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  flex-grow: 1;
}

.input-wrapper {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  margin-bottom: 40px; /* Space for command panel */
}

.input-wrapper :deep(.n-input) {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.input-wrapper :deep(.n-input__textarea-el) {
  flex-grow: 1;
  min-height: 60px;
  max-height: none;
  resize: none;
  overflow-y: auto;
}

.command-panel {
  position: absolute;
  bottom: 8px; /* Changed from 12px to 8px */
  left: 12px;
  right: 12px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 3; /* Added this line */
}

.left-controls {
  display: flex;
  gap: 8px;
}

.model-select {
  --select-height: 24px;
  height: var(--select-height);
  width: auto;
  min-width: 100px;
}

.model-select :deep(.n-base-selection) {
  height: var(--select-height);
  line-height: var(--select-height);
  font-size: 12px;
  padding: 0 4px;
  background-color: var(--n-color); /* Added this line */
  border-color: var(--n-border-color); /* Added this line */
}

.model-select :deep(.n-base-selection-label) {
  padding: 0;
}

.model-select :deep(.n-icon) {
  color: var(--n-text-color-2);
}

.model-select :deep(.n-base-suffix) {
  font-size: 12px;
}

.model-label {
  font-size: 12px;
  color: var(--n-text-color-2);
}
</style>
