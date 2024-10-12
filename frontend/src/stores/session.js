// ktt/frontend/src/stores/session.js
import { defineStore } from 'pinia';

export const useSessionStore = defineStore('session', {
  state: () => ({
    results: [],
    formValue: {
      id: 0,
      cluster_name: "",
      address: "",
      port: "",
      cmds: "",
    }
  }),
  actions: {
    setResults(newResults) {
      this.results = newResults;
    },
    addResult(newResult) {
      this.results.push(newResult);
    },
  },
});
