<template>
  <div class="tab-container">
    <div class="tab-content">
      <!-- Your main content goes here -->
    </div>
    <n-tabs
      v-model:value="value"
      type="card"
      :addable="addable"
      :closable="closable"
      tab-style="min-width: 80px;"
      @close="handleClose"
      @add="handleAdd"
      class="bottom-tabs"
    >
      <n-tab-pane v-for="panel in panels" :key="panel" :name="panel">
        {{ panel }}
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script>
import { computed, defineComponent, ref } from "vue";

export default defineComponent({
  setup() {
    const valueRef = ref(1);
    const panelsRef = ref([1]);
    const addableRef = computed(() => {
      return {
        disabled: panelsRef.value.length >= 10,
      };
    });
    const closableRef = computed(() => {
      return panelsRef.value.length > 1;
    });
    return {
      value: valueRef,
      panels: panelsRef,
      addable: addableRef,
      closable: closableRef,
      handleAdd() {
        const newValue = Math.max(...panelsRef.value) + 1;
        panelsRef.value.push(newValue);
        valueRef.value = newValue;
      },
      handleClose(name) {
        const { value: panels } = panelsRef;
        const nameIndex = panels.findIndex((panelName) => panelName === name);
        if (!~nameIndex) return;
        panels.splice(nameIndex, 1);
        if (name === valueRef.value) {
          valueRef.value = panels[Math.min(nameIndex, panels.length - 1)];
        }
      },
    };
  },
});
</script>

<style scoped>
.tab-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tab-content {
  flex-grow: 1;
  overflow-y: auto;
}

.bottom-tabs {
  flex-shrink: 0;
  height: 50px; /* Adjust this value to set the desired fixed height for the tab bar */
}
</style>
