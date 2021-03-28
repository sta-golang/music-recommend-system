import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../components/login.vue'
import Info from '../components/info.vue'

Vue.use(VueRouter)

const routes = [
  { path: '/login', component: Login },
  {
    path: '/info',
    component: Info
  }
]

const router = new VueRouter({
  routes
})

export default router
