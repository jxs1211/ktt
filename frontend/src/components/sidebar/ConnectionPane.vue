<script setup>
import useDialogStore from "stores/dialog.js";
import { NIcon, useThemeVars } from "naive-ui";
import AddGroup from "@/components/icons/AddGroup.vue";
import AddLink from "@/components/icons/AddLink.vue";
import IconButton from "@/components/common/IconButton.vue";
import Filter from "@/components/icons/Filter.vue";
import ConnectionTree from "./ConnectionTree.vue";
import { ref } from "vue";
import More from "@/components/icons/More.vue";
import Import from "@/components/icons/Import.vue";
import { useRender } from "@/utils/render.js";
import Export from "@/components/icons/Export.vue";
import useConnectionStore from "stores/connections.js";

const themeVars = useThemeVars();
const dialogStore = useDialogStore();
const connectionStore = useConnectionStore();
const render = useRender();
const filterPattern = ref("");

const moreOptions = [
  { key: "import", label: "interface.import_conn", icon: Import },
  { key: "export", label: "interface.export_conn", icon: Export },
];

const onSelectOptions = async (select) => {
  switch (select) {
    case "import":
      await connectionStore.importConnections();
      await connectionStore.initConnections(true);
      break;
    case "export":
      await connectionStore.exportConnections();
      break;
  }
};
</script>

<template>
  <div class="nav-pane-container flex-box-v">
    <connection-tree :filter-pattern="filterPattern" />
  </div>
</template>

<style lang="scss" scoped>
.nav-pane-bottom {
  color: v-bind("themeVars.iconColor");
  border-top: v-bind("themeVars.borderColor") 1px solid;
}
</style>
