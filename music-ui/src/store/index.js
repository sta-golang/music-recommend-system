import Vue from "vue";
import Vuex from "vuex";
import { _getUserSongList } from "@/network/user";
import { _getUserPlaylist } from "../network/user";
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
      console.log(obj);
      store.user = obj;
      _getUserPlaylist().then(result => {
        console.log(result);
        store.songList = result.data;
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
