// ktt/frontend/src/stores/session.js
import { defineStore } from "pinia";

export const useSessionStore = defineStore("session", {
  state: () => ({
    results: [],
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
