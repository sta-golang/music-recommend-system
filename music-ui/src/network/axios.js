import Vue from "vue";
import axios from "axios";
import { Loading } from "element-ui";
let vm = new Vue();
let loading;
let timer = null;
axios.defaults.baseURL = "http://127.0.0.1:8888/";

axios.interceptors.request.use(config => {
  loading = Loading.service({
    lock: true,
    text: "加载中……",
    background: "rgba(0, 0, 0, .3)",
    customClass: "loading"
  });
  return config;
});
axios.interceptors.response.use(
  function(res) {
    clearTimeout(timer);
    timer = setTimeout(() => {
      loading.close();
    }, 500);
    var data = res.data;
    return data;
  },
  function(error) {
    loading.close();
    // 对请求错误做些什么
    vm.$message.error("请求失败");
    return Promise.reject(error);
  }
);
Vue.prototype.$http = axios;
export default axios;
