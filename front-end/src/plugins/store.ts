import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';

Vue.use(Vuex);

interface State {
  siteKey: string | null;
  isLogin: boolean;
  progressValue: number;
  autoRefresh: boolean
}

const store: StoreOptions<State> = {
  state: {
    siteKey: null,
    isLogin: false,
    progressValue: 0,
    autoRefresh: false
  },
  getters: {
  },
  mutations: {
    setSiteKey(state: State, key: string): void {
      state.siteKey = key;
    },
    setIsLogin(state: State, bool: boolean): void {
      state.isLogin = bool;
    },
  },
  actions: {
    // setSiteKey({ commit }, key: string): void {
    //   commit('SET_SITE_KEY', key);
    // },
  },
  modules: {
  },
};

export default new Vuex.Store<State>(store);
