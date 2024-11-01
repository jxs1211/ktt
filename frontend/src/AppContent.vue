<script setup>
import ContentPane from "./components/content/ContentPane.vue";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { debounce, isEmpty } from "lodash";
import { useThemeVars } from "naive-ui";
import Ribbon from "./components/sidebar/Ribbon.vue";
import ConnectionPane from "./components/sidebar/ConnectionPane.vue";
import ContentServerPane from "./components/content/ContentServerPane.vue";
import useTabStore from "./stores/tab.js";
import usePreferencesStore from "./stores/preferences.js";
import ContentLogPane from "./components/content/ContentLogPane.vue";
import ContentErrorPane from "./components/content/ContentErrorPane.vue";
import CliBar from "@/components/content_value/CliBar.vue";
import ToolbarControlWidget from "@/components/common/ToolbarControlWidget.vue";
import {
  EventsOn,
  WindowIsFullscreen,
  WindowIsMaximised,
  WindowToggleMaximise,
} from "wailsjs/runtime/runtime.js";
import { isMacOS } from "@/utils/platform.js";
import iconUrl from "@/assets/images/icon.png";
import ResizeableWrapper from "@/components/common/ResizeableWrapper.vue";
import { extraTheme } from "@/utils/extra_theme.js";
import useConnectionStore from "./stores/connections.js";
import useConfigStore from "./stores/config.js";

const themeVars = useThemeVars();

const props = defineProps({
  loading: Boolean,
});

const data = reactive({
  navMenuWidth: 50,
  toolbarHeight: 38,
});

const tabStore = useTabStore();
const connectionStore = useConnectionStore();
const configStore = useConfigStore();
const prefStore = usePreferencesStore();
// feat: errorPaneRef support to trigger fetching data from backend using watchEffect
const errorPaneRef = ref(null);
const exThemeVars = computed(() => {
  return extraTheme(prefStore.isDark);
});
// const preferences = ref({})
// provide('preferences', preferences)

const saveSidebarWidth = debounce(prefStore.savePreferences, 1000, {
  trailing: true,
});
const handleResize = () => {
  saveSidebarWidth();
};

// watchEffect(() => {
//   if (connectionStore.clusters.length > 0) {
//     errorPaneRef.value?.refresh();
//   }
// });

const logoWrapperWidth = computed(() => {
  return `${data.navMenuWidth + prefStore.behavior.asideWidth - 4}px`;
});

const logoPaddingLeft = ref(10);
const maximised = ref(false);
const hideRadius = ref(false);
const wrapperStyle = computed(() => {
  return hideRadius.value
    ? {}
    : {
        border: `1px solid ${themeVars.value.borderColor}`,
        borderRadius: "10px",
      };
});
const spinStyle = computed(() => {
  return hideRadius.value
    ? {
        backgroundColor: themeVars.value.bodyColor,
      }
    : {
        backgroundColor: themeVars.value.bodyColor,
        borderRadius: "10px",
      };
});

const onToggleFullscreen = (fullscreen) => {
  hideRadius.value = fullscreen;
  if (fullscreen) {
    logoPaddingLeft.value = 10;
  } else {
    logoPaddingLeft.value = isMacOS() ? 70 : 10;
  }
};

const onToggleMaximize = (isMaximised) => {
  if (isMaximised) {
    maximised.value = true;
    if (!isMacOS()) {
      hideRadius.value = true;
    }
  } else {
    maximised.value = false;
    if (!isMacOS()) {
      hideRadius.value = false;
    }
  }
};

EventsOn("window_changed", (info) => {
  const { fullscreen, maximised } = info;
  onToggleFullscreen(fullscreen === true);
  onToggleMaximize(maximised);
});

onMounted(async () => {
  const fullscreen = await WindowIsFullscreen();
  onToggleFullscreen(fullscreen === true);
  const maximised = await WindowIsMaximised();
  onToggleMaximize(maximised);
  // load config from local
  const { success, msg, data } = await configStore.getLocalConfig();
  console.log("load config from local", success, msg);
  if (success) {
    connectionStore.updateClusterFromConfig(data);
  } else {
    console.error(msg);
  }
  console.log("AppContent mounted, tabs: ", tabStore.tabs);
});

const onKeyShortcut = (e) => {
  switch (e.key) {
    case "w":
      if (e.metaKey) {
        // close current tab
        const tabStore = useTabStore();
        const currentTab = tabStore.currentTab;
        if (currentTab != null) {
          tabStore.closeTab(currentTab.name);
        }
      }
      break;
  }
};
watch(
  () => tabStore.nav,
  (newVal, oldVal) => {
    console.log("watch nav: ", tabStore.currentTab, tabStore.currentSubTab);
    if (newVal === "browser") {
      tabStore.currentSubTab = "status";
    }
  },
);
</script>

<template>
  <!-- app content-->
  <n-spin
    :show="props.loading"
    :style="spinStyle"
    :theme-overrides="{ opacitySpinning: 0 }"
  >
    <div
      id="app-content-wrapper"
      :style="wrapperStyle"
      class="flex-box-v"
      tabindex="0"
      @keydown="onKeyShortcut"
    >
      <!-- title bar -->
      <div
        id="app-toolbar"
        :style="{ height: data.toolbarHeight + 'px' }"
        class="flex-box-h"
        style="--wails-draggable: drag"
        @dblclick="WindowToggleMaximise"
      >
        <!-- title -->
        <div
          id="app-toolbar-title"
          :style="{
            width: logoWrapperWidth,
            minWidth: logoWrapperWidth,
            paddingLeft: `${logoPaddingLeft}px`,
          }"
        >
          <n-space :size="3" :wrap="false" :wrap-item="false" align="center">
            <n-avatar
              :size="28"
              :src="iconUrl"
              color="#0000"
              style="min-width: 28px"
            />
            <div style="min-width: 30px; white-space: nowrap; font-weight: 800">
              KT
            </div>
            <n-gradient-text
              type="success"
              :size="16"
              v-if="tabStore.nav !== 'log' && !isEmpty(tabStore.currentTabName)"
            >
              - {{ tabStore.currentTabName }}
            </n-gradient-text>
          </n-space>
        </div>
        <div class="flex-item-expand" style="min-width: 15px"></div>
        <!-- simulate window control buttons -->
        <toolbar-control-widget
          v-if="!isMacOS()"
          :maximised="maximised"
          :size="data.toolbarHeight"
          style="align-self: flex-start"
        />
      </div>

      <!-- content -->
      <div
        id="app-content"
        :style="prefStore.generalFont"
        class="flex-box-h flex-item-expand"
        style="--wails-draggable: none"
      >
        <ribbon v-model:value="tabStore.nav" :width="data.navMenuWidth" />

        <!-- browser -->
        <div
          v-show="tabStore.nav === 'browser'"
          class="content-area flex-box-h flex-item-expand"
        >
          <!-- <content-value-tab /> -->
          <content-pane
            v-for="t in tabStore.tabs"
            v-show="tabStore.currentTabName === t.name"
            :key="t.name"
            :server="t.name"
            class="flex-item-expand"
          />
          <!-- Bottom Bar -->
          <div class="bottom-bar">
            <!-- You can add content here, like buttons or text -->
            <span>Bottom Bar Content</span>
          </div>
        </div>

        <!-- server list page -->
        <div
          v-show="tabStore.nav === 'server'"
          class="content-area flex-box-h flex-item-expand"
        >
          <resizeable-wrapper
            v-model:size="prefStore.behavior.asideWidth"
            :min-size="250"
            :offset="data.navMenuWidth"
            class="flex-item"
          >
            <connection-pane class="app-side flex-item-expand" />
          </resizeable-wrapper>
          <content-server-pane class="flex-item-expand" />
        </div>

        <!-- cli -->
        <div
          v-show="tabStore.nav === 'cli'"
          class="content-area flex-box-h flex-item-expand"
        >
          <CliBar :min-height="'50px'" />
          <!-- <resizeable-wrapper
            v-model:size="prefStore.behavior.asideWidth"
            :min-size="250"
            :offset="data.navMenuWidth"
            class="flex-item"
          >
          </resizeable-wrapper> -->
        </div>

        <!-- log -->
        <!-- <div
          v-show="tabStore.nav === 'log'"
          class="content-area flex-box-h flex-item-expand"
        >
          <content-log-pane ref="logPaneRef" class="flex-item-expand" />
        </div> -->
      </div>
    </div>
  </n-spin>
</template>

<style lang="scss" scoped>
#app-content-wrapper {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  box-sizing: border-box;
  background-color: v-bind("themeVars.bodyColor");
  color: v-bind("themeVars.textColorBase");

  #app-toolbar {
    background-color: v-bind("exThemeVars.titleColor");
    border-bottom: 1px solid v-bind("exThemeVars.splitColor");

    &-title {
      padding-left: 10px;
      padding-right: 10px;
      box-sizing: border-box;
      align-self: center;
      align-items: baseline;
    }
  }

  .app-toolbar-tab {
    align-self: flex-end;
    margin-bottom: -1px;
    margin-left: 3px;
    overflow: auto;
  }

  #app-content {
    height: calc(100% - 60px);

    .content-area {
      overflow: hidden;
    }
  }

  .app-side {
    //overflow: hidden;
    height: 100%;
    background-color: v-bind("exThemeVars.sidebarColor");
    border-right: 1px solid v-bind("exThemeVars.splitColor");
  }
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.bottom-bar {
  height: 20px; /* Set a fixed height for the bar */
  background-color: #f0f0f0; /* Change to your desired color */
  border-top: 1px solid #ccc; /* Optional: Add a border for separation */
  display: flex;
  align-items: center; /* Center content vertically */
  justify-content: center; /* Center content horizontally */
}
</style>
