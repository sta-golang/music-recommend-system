<template>
  <div>
    <el-container>
      <!-- 头部信息 -->
      <el-header>
        <header-temp></header-temp>
      </el-header>
      <el-container>
        <!-- 侧边栏 -->
        <!-- style="overflow-y:hidden;" -->
        <el-aside width="200px">
          <aside-menu> </aside-menu>
        </el-aside>
        <!-- 主体部分 -->
        <el-main>
          <!-- 通过路由匹配切换组件 -->
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
    <musicplay></musicplay>
  </div>
</template>

<script>
// 侧边栏组件
import asideMenu from './aside/Aside.vue'
// 头部
import headerTemp from './header/header'
// 音乐播放器
import Musicplay from './centent/musicplay/Musicplay'
export default {

  components: {
    asideMenu,
    headerTemp,
    Musicplay
  }
}
</script>

<style lang='less'>
.el-header {
  z-index: 150;
  padding: 0 !important;
}
.el-main {
  margin-bottom: 90px;
}
</style>
