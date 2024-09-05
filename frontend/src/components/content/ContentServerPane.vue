<script setup>
import { isEmpty } from "lodash";
import AddLink from "@/components/icons/AddLink.vue";
import useDialogStore from "@/stores/dialog.js";
import useConnectionStore from "@/stores/connections.js";

const dialogStore = useDialogStore();
const connectionStore = useConnectionStore();
</script>

<template>
  <div class="content-container flex-box-v">
    <!-- TODO: replace icon to app icon -->
    <template v-if="isEmpty(connectionStore.clusters)">
      <n-empty :description="$t('interface.empty_server_content')">
        <template #extra>
          <n-button :focusable="false" @click="dialogStore.openNewConfigDialog()">
            <template #icon>
              <n-icon :component="AddLink" size="18" />
            </template>
            {{ $t("interface.new_conn") }}
          </n-button>
        </template>
      </n-empty>
    </template>
    <template v-else>
      <template v-if="isEmpty(connectionStore.currentCluster)">
        <n-empty :description="$t('interface.choose_one_conn')">
        </n-empty>
      </template>
      <template v-else-if="!connectionStore.switchedClusterOK">
        <n-spin class="spinner">
          <!-- Add your custom spinner here -->
        </n-spin>
      </template>
    </template>
  </div>
</template>

<style lang="scss" scoped>
@import "@/styles/content";
.spinner {
  height: 100%;
  width: 100%;
  position: relative;
}
.content-container {
  justify-content: center;
  padding: 5px;
  box-sizing: border-box;
}
</style>
