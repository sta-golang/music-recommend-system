import Vue from 'vue'
import App from './App.vue'
import router from './router'
import './plugins/element.js'
import axios from 'axios'

axios.defaults.baseURL = 'http://localhost:8888'
Vue.prototype.$http = axios
Vue.prototype.$tokenStr = 'sta-token'

axios.defaults.headers.common[Vue.prototype.$tokenStr] = localStorage.getItem(
  Vue.prototype.$tokenStr
)
Vue.config.productionTip = false
new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
