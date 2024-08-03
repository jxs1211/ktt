<template>
  <div v-for="(item, key) in data.solution" :key="key">
    <div class="item-container">
      <span class="item-content">
        {{ item }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from "vue";
const props = defineProps({
  data: {
    type: Object,
    default: () => ({
      // error: 'Error: --- Back-off pulling image "nginx:latest" ---',
      solution: [],
    }),
  },
});
</script>

<style scoped>
.item-container {
  overflow-x: auto;
  text-align: left;
  padding: 2px;
}
.item-content {
  display: inline-block;
  white-space: nowrap;
}
</style>
