import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store/index";
import "./plugins/element.js";
import "@/network/axios";
import "./common/css/base.css";
import "./common/css/global.css";
import "default-passive-events";
import axios from "axios";
Vue.config.productionTip = false;
Vue.prototype.$bus = new Vue();
Vue.prototype.$tokenStr = "sta-token";
Vue.prototype.$globalLoginState = false;
Vue.prototype.$http = axios;
axios.defaults.headers.common[Vue.prototype.$tokenStr] = store.state.token;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
