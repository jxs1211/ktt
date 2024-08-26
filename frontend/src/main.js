import { createPinia } from "pinia";
import { createApp, nextTick } from "vue";
import App from "./App.vue";
import "./styles/style.scss";
import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import relativeTime from "dayjs/plugin/relativeTime";
import { i18n } from "@/utils/i18n.js";
import { setupDiscreteApi } from "@/utils/discrete.js";
import usePreferencesStore from "stores/preferences.js";
import { loadEnvironment } from "@/utils/platform.js";
import { setupMonaco } from "@/utils/monaco.js";
import { setupChart } from "@/utils/chart.js";
import useConfigStore from "./stores/config";

dayjs.extend(duration);
dayjs.extend(relativeTime);

async function setupApp() {
  const app = createApp(App);
  app.use(i18n);
  app.use(createPinia());

  await loadEnvironment();
  setupMonaco();
  setupChart();
  const prefStore = usePreferencesStore();
  await prefStore.loadPreferences();
  await setupDiscreteApi();
  const configStore = useConfigStore();
  const { success, msg, _ } = await configStore.loadConfigFromLocal();
  if (!success) {
    console.error("load config from local failed: ", msg);
  } else {
    console.log("load config from local success");
  }
  app.config.errorHandler = (err, instance, info) => {
    // TODO: add "send error message to author" later
    nextTick().then(() => {
      try {
        const content = err.toString();
        $notification.error(content, {
          title: i18n.global.t("common.error"),
          meta: "Please see console output for more detail",
        });
        console.error(err);
      } catch (e) {}
    });
  };
  // app.config.warnHandler = (message) => {
  //     console.warn(message)
  // }
  app.mount("#app");
}

setupApp();
