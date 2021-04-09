import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../components/Home.vue";
import Discover from "../components/discover/Discover.vue";
import Allmv from "../components/allmv/Allmv.vue";

// 详情歌单页面
import Musiclistdetail from "@/views/musiclistdetail/Musiclistdetail";
// 歌手详情
import ArtistDetail from "@/views/artistdetail/ArtistDetail";
// 搜索结果
import SearchList from "@/views/searchdetail/SearchList";
Vue.use(VueRouter);

const routes = [
  { path: "/", redirect: "/home/discover" },
  {
    path: "/home",
    component: Home,
    children: [
      { path: "/home/discover", component: Discover },
      { path: "/home/allmv", component: Allmv },
      { path: "/home/musiclistdetail", component: Musiclistdetail },
      { path: "/home/artistalbum", component: ArtistDetail },
      { path: "/home/searchlist", component: SearchList }
    ]
  }
];
const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err);
};
const router = new VueRouter({
  routes
});

export default router;
