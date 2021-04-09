<template>
  <div class="aside">
    <dl>
      <dt>推荐</dt>
      <dd :class="['cursorPointer',current=='discover' ? 'current':'']" @click="toPath('discover')"><i class="iconfont">&#xe680;</i>发现音乐</dd>
      <dd :class="['cursorPointer',current=='allmv' ? 'current':'']" @click="toPath('allmv')"><i class="iconfont">&#xe6c4;</i>全部MV</dd>
      <!-- <dd class="cursorPointer" @click="toPath('allmv')"><i class="iconfont">&#xe670;</i>直播</dd> -->
      <!-- <dd class="cursorPointer" @click="toPath('allmv')"><i class="iconfont">&#xe614;</i>视频</dd> -->
      <!-- <dd class="cursorPointer" @click="toPath('allmv')"><i class="iconfont">&#xe61a;</i>朋友</dd> -->
    </dl>
    <!-- <dl>
      <dt>我的音乐</dt>
      <dd class="cursorPointer" @click="toPath('allmv')"><i class="iconfont">&#xe601;</i>本地音乐</dd>
      <dd class="cursorPointer" @click="toPath('allmv')"><i class="iconfont">&#xe723;</i>下载管理</dd>
    </dl> -->
    <dl>
      <dt>我的歌单</dt>
      <dd :class="['cursorPointer',current==item.id ? 'current':'']" v-for="item in songList" :key="item.id" @click="toSong(item.id)"><i class="iconfont">&#xe83e;</i>{{item.name}}</dd>
    </dl>
  </div>
</template>

<script>
export default {
  data () {
    return {
      current: 'discover',

    }
  },
  methods: {
    toPath (path) {
      // window.sessionStorage.setItem('path', path)
      this.current = path
      this.$router.push(path)
    },
    toSong (id) {
      this.current = id
      this.$router.push({ path: '/home/musiclistdetail', query: { id } })
    }
  },
  created () {
    // this.current = window.sessionStorage.getItem('path')
  },
  computed: {
    songList () {
      return this.$store.state.songList
    }
  }
}
</script>

<style lang='less' scoped>
.aside {
  margin-bottom: 160px;
  dl {
    margin: 0;
    margin: 10px 0;
    dt {
      padding-left: 10px;
      height: 33px;
      font-size: 14px;
    }
    dd {
      width: 200px;
      box-sizing: border-box;
      text-overflow: ellipsis;
      white-space: nowrap;
      overflow: hidden;
      color: #555;
      height: 33px;
      line-height: 33px;
      font-size: 15px;
      padding-left: 18px;
      border-left: 2px solid transparent;
      transition: 0.2s all;
      i {
        font-size: 18px;
        margin-right: 10px;
      }
      &:hover {
        color: #000;
      }
    }
    dd.current {
      color: #222;
      border-left: 2px solid #c62f2f;
      background-color: #e6e7ea;
    }
  }
}
</style>
