import Vue from 'vue'
import {
  Button,
  Header,
  Container,
  Aside,
  Main,
  Input,
  Dialog,
  Form,
  FormItem,
  Message,
  Avatar,
  Tabs,
  TabPane,
  Table,
  TableColumn,
  Slider,
  Drawer,
  Tooltip
} from 'element-ui'
import infiniteScroll from "vue-infinite-scroll";
Vue.use(infiniteScroll)
Vue.use(Drawer)
Vue.use(Tooltip)
Vue.use(TabPane)
Vue.use(Slider)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Tabs)
Vue.use(Avatar)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Dialog)
Vue.use(Button)
Vue.use(Input)
Vue.use(Header)
Vue.use(Container)
Vue.use(Aside)
Vue.use(Main)

Vue.prototype.$message = Message
