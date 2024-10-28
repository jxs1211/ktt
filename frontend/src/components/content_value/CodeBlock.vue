<template>
  <pre><code :class="language" ref="codeBlock">{{ code }}</code></pre>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import hljs from "highlight.js";
import "highlight.js/styles/vs2015.css"; // Choose a style that fits your design

const props = defineProps({
  code: {
    type: String,
    required: true,
  },
  language: {
    type: String,
    default: "plaintext",
  },
});

const codeBlock = ref(null);

const highlightCode = () => {
  if (codeBlock.value) {
    hljs.highlightElement(codeBlock.value);
  }
};

onMounted(() => {
  highlightCode();
});

watch(
  () => props.code,
  () => {
    nextTick(highlightCode);
  },
);
</script>

<style scoped>
pre {
  margin: 0;
  padding: 10px;
  background-color: #1e1e1e; /* Dark background for code */
  border-radius: 4px;
  overflow-x: auto;
}

code {
  font-family: "Courier New", Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
}
</style>
