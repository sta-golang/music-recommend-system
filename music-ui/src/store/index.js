import Vue from "vue";
import Vuex from "vuex";
import { _getUserSongList } from "@/network/user";
Vue.use(Vuex);
export default new Vuex.Store({
  // 在里面存放的就是全局共享的数据
  state: {
    user: {},
    cookie: "",
    uId: "",
    songList: [],
    showLyric: false,
    songDetail: {},
    isMusicPlay: false,
    currentTime: 0,
    token: ""
  },
  mutations: {
    addUser(store, obj) {
      store.user = obj;
      store.cookie = store.user.cookie;
      store.uId = store.user.account.id;
      _getUserSongList(store.uId).then(result => {
        store.songList = result.playlist;
      });
    },
    storeToken(store, token) {
      store.token = token;
    },
    editshowLyric(store, type) {
      store.showLyric = type;
    },
    editSongDetai(store, obj) {
      store.songDetail = obj;
    },
    editMusicPlay(store, type) {
      store.isMusicPlay = type;
    },
    editCurrentTime(store, time) {
      store.currentTime = time;
    }
  }
});
