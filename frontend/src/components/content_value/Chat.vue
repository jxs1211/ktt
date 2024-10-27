<template>
  <div class="chat-container">
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
    <div class="message-display flex-item" ref="messageDisplayRef">
      <div v-for="message in chatMessages" :key="message.id" :class="['message', message.type]">
        <!-- <div class="message-content" v-if="message.content.trim()">
          {{ message.content }}
        </div> -->
        <template v-if="message.content.trim()">
          {{ message.content }}
        </template>
        <n-button text size="tiny" @click="copyMessage(message.content)" class="copy-button">
          Copy
        </n-button>
      </div>
    </div>
    <div class="chat-input-container">
      <n-input
        ref="inputRef"
        v-model:value="inputMessage"
        type="textarea"
        placeholder="Add context"
        :disabled="isWaiting"
        @keydown.enter.prevent="handleEnterKey"
        @input="adjustTextareaHeight"
        class="chat-input"
      />
      <div class="command-panel">
        <n-select
          v-model:value="selectedModel"
          :options="modelOptions"
          size="tiny"
          filterable
          clearable
          :render-label="renderModelLabel"
          class="model-select"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { EventsOff, EventsOn } from 'wailsjs/runtime/runtime.js';
import { watch, ref, h, nextTick, onMounted, onUnmounted, computed } from "vue";
import { NButton, NIcon, NInput, NSelect } from "naive-ui";
import { CaretUp, Add, Time, Close, ChevronDown } from "@vicons/ionicons5";
import CodeBlock from "./CodeBlock.vue";
import { useThemeVars } from "naive-ui";
import { GetCompletion2 } from "wailsjs/go/ai/ClientService.js";
import usePreferencesStore from '../../stores/preferences';
// import hljs from 'highlight.js';

// hljs.configure({
//   languages: ['go', 'javascript', 'python', 'bash', 'json'], // add the languages you need
// });

/**
 * Resizeable component wrapper
 */
const themeVars = useThemeVars();
const messages = ref([]);
const inputMessage = ref("");
const chatMessages = ref([])  // New ref to hold all chat messages
const isWaiting = ref(false);  // New ref to track waiting state
const props = defineProps({
  session: {
    type: String,
    required: true,
  }
});
// Model selection
const selectedModel = ref("llama3.2");
const supportedModelsOptions = computed(() => {

})
const modelOptions = [
  {
    label: "llama3.2",
    value: "llama3.2",
    icon: "cursor-icon", // You might want to replace this with actual icon
  },
  // {
  //   label: "cursor-medium",
  //   value: "cursor-medium",
  //   icon: "cursor-icon",
  // },
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
  if (inputRef.value) {
    const textarea = inputRef.value.$el.querySelector("textarea");
    textarea.style.height = 'auto';
    textarea.style.height = `${Math.min(Math.max(textarea.scrollHeight, 60), 200)}px`;
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

const sendMessage = async () => {
  if (inputMessage.value.trim() === "") {
    console.log("input msg is empty")
    return;
  }
  // Add user message to chat
  chatMessages.value.push({
    id: Date.now(),
    type: 'user',
    content: inputMessage.value.trim()
  });

  const userMessage = inputMessage.value.trim();
  inputMessage.value = "";

  // Show waiting message
  isWaiting.value = true;

  try {
    // Send message to backend
    await GetCompletion2(props.session, selectedModel.value, userMessage);
    // Note: We don't handle the response here anymore
    // The response will be handled by the event listener
  } catch (error) {
    console.error("Get completion failed: ", error);
    // Remove waiting message
    chatMessages.value = chatMessages.value.filter(msg => msg.type !== 'system');
    // Add error message to chat
    chatMessages.value.push({
      id: Date.now() + 2,
      type: 'system',
      content: 'Error: Unable to get response'
    });
  } finally {
    isWaiting.value = false;
  }

  nextTick(() => {
    adjustTextareaHeight();
  });
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
const eventName = `ai:chat:${props.session}`;
const messageDisplayRef = ref(null);

const scrollToBottom = () => {
  nextTick(() => {
    if (messageDisplayRef.value) {
      messageDisplayRef.value.scrollTop = messageDisplayRef.value.scrollHeight;
    }
  });
};

// Watch for changes in chatMessages and scroll to bottom
watch(() => chatMessages.value.length, scrollToBottom);

onMounted(() => {
  adjustTextareaHeight();
  EventsOn(eventName, (data) => {
    console.log("on " + eventName + " : ", data);
    // Remove waiting message
    chatMessages.value = chatMessages.value.filter(msg => msg.type !== 'system');

    if (typeof data === 'string' && data.startsWith('systemError:')) {
      // Handle error case
      chatMessages.value.push({
        id: Date.now(),
        type: 'system',
        content: data
      });
    } else {
      // Handle success case
      chatMessages.value.push({
        id: Date.now(),
        type: 'ai',
        content: data
      });
    }
    isWaiting.value = false;
    scrollToBottom();
  });
});
onUnmounted(() => {
  EventsOff(eventName);
});
// Watch isWaiting for changes
watch(isWaiting, (newValue) => {
  console.log('isWaiting changed:', newValue);
});

const copyMessage = (content) => {
  navigator.clipboard.writeText(content).then(() => {
    // You can use a toast notification here if you have one
    console.log('Message copied to clipboard');
  }).catch(err => {
    console.error('Failed to copy message: ', err);
  });
};

</script>

<style scoped>
.chat-container {
  --chat-background-color: #c3baba00; /* Define the variable here */
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--chat-background-color); /* Use the variable */
  border-left: 1px solid var(--n-border-color);
  box-sizing: border-box;
  overflow: hidden;
  position: relative;
  user-select: text; /* Explicitly allow text selection */
}

.chat-header {
  flex-shrink: 0;
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
  flex-grow: 1;
  overflow-y: auto;
  padding: 16px;
  padding-bottom: 120px; /* Increased to give more space for the input */
  width: 100%;
  box-sizing: border-box;
  cursor: text; /* Show text cursor to indicate selectability */
}

.chat-input-container {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px;
  background-color: var(--chat-background-color);
  border-top: 1px solid var(--n-border-color);
  display: flex;
  flex-direction: column;
}

.chat-input {
  width: 100%;
  margin-bottom: 8px;
}

.chat-input :deep(.n-input__textarea-el) {
  min-height: 60px;
  max-height: 200px;
  resize: none;
}

.command-panel {
  display: flex;
  justify-content: flex-start;
  align-items: center;
}

.model-select {
  width: 200px;
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

.system {
  font-style: italic;
  color: #888;
}

.user {
  text-align: left;
  background-color: rgba(0, 0, 0, 0.05);
  padding: 8px;
  margin-bottom: 8px;
}

.ai {
  text-align: left;
  background-color: rgba(0, 0, 255, 0.05);
  padding: 8px;
  margin-bottom: 8px;
}

/* Ensure the textarea doesn't overflow */
.input-wrapper :deep(.n-input__textarea-el) {
  max-height: 300px; /* Adjust this value as needed */
  overflow-y: auto;
}

.message-display :deep(pre) {
  background-color: #f4f4f4;
  padding: 1em;
  border-radius: 5px;
  overflow-x: auto;
}

.message-display :deep(code) {
  background-color: #f4f4f4;
  padding: 0.2em 0.4em;
  border-radius: 3px;
}

.message-display :deep(table) {
  border-collapse: collapse;
  margin: 1em 0;
}

.message-display :deep(th), .message-display :deep(td) {
  border: 1px solid #ddd;
  padding: 8px;
}

.message-display :deep(blockquote) {
  border-left: 4px solid #ddd;
  padding-left: 1em;
  color: #666;
}

.message-display > div:empty {
  display: none;
}

.message {
  margin-bottom: 8px;
  padding: 8px;
  border-radius: 4px;
}

.message-content {
  white-space: pre-wrap; /* Preserve whitespace and allow wrapping */
  word-break: break-word; /* Break long words to prevent overflow */
}

/* Style for selected text */
::selection {
  background-color: rgba(0, 123, 255, 0.3); /* Light blue background for selected text */
  color: inherit; /* Keep the text color */
}

.message {
  position: relative;
}

.copy-button {
  position: absolute;
  top: 4px;
  right: 4px;
  opacity: 0.3; /* Partially visible by default */
  transition: opacity 0.2s;
}

.message:hover .copy-button {
  opacity: 1; /* Fully visible on hover */
}
</style>
