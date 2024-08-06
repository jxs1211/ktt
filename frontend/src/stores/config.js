import { defineStore } from "pinia";
import { get, isEmpty, isObject, uniq } from "lodash";
import {
  LoadConfig,
  GetLocalConfig,
  TestConnection,
  Analyze,
  GetClusters,
  GetAvailableFilteredResources,
} from "wailsjs/go/client/ClientService.js";
import { ConnectionType } from "@/consts/connection_type.js";
import useBrowserStore from "stores/browser.js";
import { i18nGlobal } from "@/utils/i18n.js";
import { ClipboardGetText } from "wailsjs/runtime/runtime.js";

const useConfigStore = defineStore("config", {
  state: () => ({
    loaded: false,
    testResult: "",
    localConfig: {},
  }),
  getters: {},
  actions: {
    async getAvailableResources() {
      return await GetAvailableFilteredResources();
    },
    async getLocalConfig() {
      return await GetLocalConfig();
    },
    async loadConfig(content) {
      if (isEmpty(content)) {
        return { success: false, msg: "no config content" };
      }
      return await LoadConfig(content);
    },
    async testConnection(config) {
      if (isEmpty(config)) {
        return { success: false, msg: "no config content" };
      }
      const { success, msg, data } = await TestConnection(config);
      return { success, msg, data };
    },
    async analyze(
      cluster,
      aibackend,
      model,
      fitlers,
      explain,
      aggregate,
      anonymize,
    ) {
      return await Analyze(
        cluster,
        aibackend,
        model,
        fitlers,
        explain,
        aggregate,
        anonymize,
      );
    },
    async getClusters() {
      return await GetClusters();
    },
    /**
     * get connection by name from local profile
     * @param name
     * @returns {Promise<ConnectionProfile|null>}
     */
    async getConnectionProfile(name) {
      try {
        const { data, success } = await GetConnection(name);
        if (success) {
          this.serverProfile[name] = {
            defaultFilter: data.defaultFilter,
            keySeparator: data.keySeparator,
            markColor: data.markColor,
          };
          return data;
        }
      } finally {
      }
      return null;
    },

    /**
     * create a new connection or update current connection profile
     * @param {string} name set null if create a new connection
     * @param {{}} param
     * @returns {Promise<{success: boolean, [msg]: string}>}
     */
    async saveConnection(name, param) {
      const { success, msg } = await SaveConnection(name, param);
      if (!success) {
        return { success: false, msg };
      }

      // reload connection list
      await this.initConnections(true);
      return { success: true };
    },

    /**
     * save connection after sort
     * @returns {Promise<void>}
     */
    async saveConnectionSorted() {
      const mapToList = (conns) => {
        const list = [];
        for (const conn of conns) {
          if (conn.type === ConnectionType.Group) {
            const children = mapToList(conn.children);
            list.push({
              name: conn.label,
              type: "group",
              connections: children,
            });
          } else if (conn.type === ConnectionType.Server) {
            list.push({
              name: conn.name,
            });
          }
        }
        return list;
      };
      const s = mapToList(this.connections);
      SaveSortedConnection(s);
    },

    /**
     * remove connection
     * @param name
     * @returns {Promise<{success: boolean, [msg]: string}>}
     */
    async deleteConnection(name) {
      // close connection first
      const browser = useBrowserStore();
      await browser.closeConnection(name);
      const { success, msg } = await DeleteConnection(name);
      if (!success) {
        return { success: false, msg };
      }
      await this.initConnections(true);
      return { success: true };
    },

    /**
     * create a connection group
     * @param name
     * @returns {Promise<{success: boolean, [msg]: string}>}
     */
    async createGroup(name) {
      const { success, msg } = await CreateGroup(name);
      if (!success) {
        return { success: false, msg };
      }
      await this.initConnections(true);
      return { success: true };
    },

    /**
     * rename connection group
     * @param name
     * @param newName
     * @returns {Promise<{success: boolean, [msg]: string}>}
     */
    async renameGroup(name, newName) {
      if (name === newName) {
        return { success: true };
      }
      const { success, msg } = await RenameGroup(name, newName);
      if (!success) {
        return { success: false, msg };
      }
      await this.initConnections(true);
      return { success: true };
    },

    /**
     * delete group by name
     * @param {string} name
     * @param {boolean} [includeConn]
     * @returns {Promise<{success: boolean, [msg]: string}>}
     */
    async deleteGroup(name, includeConn) {
      const { success, msg } = await DeleteGroup(name, includeConn === true);
      if (!success) {
        return { success: false, msg };
      }
      await this.initConnections(true);
      return { success: true };
    },

    /**
     * save last selected database
     * @param {string} name
     * @param {number} db
     * @return {Promise<{success: boolean, [msg]: string}>}
     */
    async saveLastDB(name, db) {
      const { success, msg } = await SaveLastDB(name, db);
      if (!success) {
        return { success: false, msg };
      }
      return { success: true };
    },

    /**
     * get default key filter pattern by server name
     * @param name
     * @return {string}
     */
    getDefaultKeyFilter(name) {
      const { defaultFilter = "*" } = this.serverProfile[name] || {};
      return defaultFilter;
    },

    /**
     * get default key separator by server name
     * @param name
     * @return {string}
     */
    getDefaultSeparator(name) {
      const { keySeparator = ":" } = this.serverProfile[name] || {};
      return keySeparator;
    },

    /**
     * get default status refresh interval by server name
     * @param {string} name
     * @return {number}
     */
    getRefreshInterval(name) {
      const { refreshInterval = 5 } = this.serverProfile[name] || {};
      return refreshInterval;
    },

    /**
     * set and save default refresh interval
     * @param {string} name
     * @param {number} interval
     * @return {Promise<{success: boolean}|{msg: undefined, success: boolean}>}
     */
    async saveRefreshInterval(name, interval) {
      const profile = this.serverProfile[name] || {};
      profile.refreshInterval = interval;
      const { success, msg } = await SaveRefreshInterval(name, interval);
      if (!success) {
        return { success: false, msg };
      }
      return { success: true };
    },

    /**
     * export connections to zip
     * @return {Promise<void>}
     */
    async exportConnections() {
      const {
        success,
        msg,
        data: { path = "" },
      } = await ExportConnections();
      if (!success) {
        if (!isEmpty(msg)) {
          $message.error(msg);
          return;
        }
      }

      $message.success(i18nGlobal.t("dialogue.handle_succ"));
    },

    /**
     * import connections from zip
     * @return {Promise<void>}
     */
    async importConnections() {
      const { success, msg } = await ImportConnections();
      if (!success) {
        if (!isEmpty(msg)) {
          $message.error(msg);
          return;
        }
      }

      $message.success(i18nGlobal.t("dialogue.handle_succ"));
    },

    /**
     * parse redis url from text in clipboard
     * @return {Promise<{}>}
     */
    async parseUrlFromClipboard() {
      const urlString = await ClipboardGetText();
      if (isEmpty(urlString)) {
        throw new Error("no text in clipboard");
      }

      const { success, msg, data } = await ParseConnectURL(urlString);
      if (!success || !isObject(data)) {
        throw new Error(msg || "unknown");
      }

      data.url = urlString;
      return data;
    },
  },
});

export default useConfigStore;
