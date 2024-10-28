// ktt/frontend/src/stores/session.js
import { defineStore } from "pinia";

export const useSessionStore = defineStore("session", {
  state: () => ({
    results: [],
    debugRowMsg: "",
    formValue: {
      cluster_name: "",
      address: "",
      port: "",
      cmds: "",
    },
    data: [],
  }),
  actions: {
    // filter from data by address, port, cmds
    filterData(address, port, cmds) {
      return data.filter((item) => {
        const addressMatch = item.address === address;
        const portMatch = item.port === port;
        const cmdsMatch = cmds.every((cmd) => item.cmds.includes(cmd));
        return addressMatch && portMatch && cmdsMatch;
      });
    },
    saveData(address, port, cmds) {},
    setResults(newResults) {
      this.results = newResults;
    },
    addResult(newResult) {
      this.results.push(newResult);
    },
    emptyResults() {
      this.results = [];
    },
  },
});
