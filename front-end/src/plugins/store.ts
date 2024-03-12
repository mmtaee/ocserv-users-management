import { User } from "@/utils/types";
import Vue from "vue";
import Vuex, { StoreOptions } from "vuex";

Vue.use(Vuex);

interface Snackbar {
  text: string | null;
  color: string | null;
}

interface State {
  user: User;
  siteKey: string | null;
  isLogin: boolean;
  progressValue: number;
  autoRefresh: boolean;
  snacbar: Snackbar;
  loadingOverlay: boolean;
  loadingOverlayText: string | null;
}

const store: StoreOptions<State> = {
  state: {
    user: {
      username: null,
      is_admin: false,
    },
    siteKey: null,
    isLogin: false,
    progressValue: 0,
    autoRefresh: false,
    snacbar: {
      text: null,
      color: null,
    },
    loadingOverlay: false,
    loadingOverlayText: null,
  },
  getters: {
    getUserstate(state): User {
      return state.user;
    },
    getGoogleSiteKey(state): string | null {
      return state.siteKey;
    },
  },
  mutations: {
    setUser(state: State, values: User) {
      state.user = Object.assign({}, values);
    },
    setSiteKey(state: State, key: string): void {
      state.siteKey = key;
    },
    setIsLogin(state: State, bool: boolean): void {
      state.isLogin = bool;
    },
    setSnackBar(state: State, data: Snackbar): void {
      state.snacbar = Object.assign({}, data);
    },
    setLoadingOverlay(
      state: State,
      data: {
        active: boolean;
        text: string;
      }
    ): void {
      state.loadingOverlay = data.active;
      state.loadingOverlayText = data.text;
    },
  },
  actions: {
    // setSiteKey({ commit }, key: string): void {
    //   commit('SET_SITE_KEY', key);
    // },
  },
  modules: {},
};

export default new Vuex.Store<State>(store);
