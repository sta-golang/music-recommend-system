<template>
  <div class="MusicPlayList" style="overflow-y:auto;">
    <div class="head">
      <div class="title">
        <div>
          <span :class="[current=='playlist'?'current':'']" @click="current='playlist'">播放列表</span>
          <span :class="[current=='historical'?'current':'']" @click="current='historical'">历史记录</span>
        </div>
        <i class="close iconfont" @click="closeHandle">&#xe66a;</i>
      </div>
      <div class="warp">
        <div class="sum">共{{sum}}首</div>
        <div class="clearAll" @click="clearAll()"><i class="iconfont">&#xe606;</i> 清空</div>
      </div>
    </div>

    <div class="playlist" v-if="current=='playlist'">

      <table-list @rowClick='rowClick' :hiddenOperation='false' :hiddenIndex='false' :tracklist='musicList' :hiddenTableHeader='false' :hiddenAlbum='false' :hiddenSinger='true' :currentIndex='currentIndex'></table-list>
    </div>
    <div class="historical" v-else>历史记录</div>
  </div>
</template>

<script>
import TableList from '../musiclist/TableList'
export default {
  props: ["currentIndex", 'musicList'],
  data () {
    return {
      // 当前选中的
      current: 'playlist',
    }
  },
  methods: {
    closeHandle () {
      this.$parent.isShowMusicplaylist = false
    },
    async clearAll () {
      // 调用父组件的渐变声音
      await this.$parent.musicGradients('down')
      this.$parent.$refs.audio.pause()
      await this.$parent.resetInfo()
    },
    rowClick (index, list) {
      this.$bus.$emit('playMusic', index, list)
    }
  },
  computed: {
    sum () {
      return this.musicList.length
    }
  },
  components: {
    TableList
  }
}
</script>

<style lang='less' scoped>
.MusicPlayList {
  position: absolute;
  bottom: 50px;
  right: 0;
  width: 300px;
  z-index: 129;
  height: 400px;
  background-color: #fafafa;
  box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.5);
  .head {
    .title {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 40px;
      position: relative;
      background-color: #f3f3f5;
      border-bottom: 1px solid #ccc;
      div {
        width: 160px;
        border-radius: 5px;
        overflow: hidden;
        border: 1px solid #7c7d85;
        span {
          display: inline-block;
          width: 80px;
          height: 25px;
          color: #7c7d85;
          background-color: #ffffff;
          text-align: center;
          line-height: 25px;
          box-sizing: border-box;
          cursor: pointer;
        }
        .current {
          color: #fff;
          background-color: #7c7d85;
        }
      }
      .close {
        position: absolute;
        top: 5px;
        right: 5px;
        cursor: pointer;
      }
    }
    .warp {
      height: 30px;
      display: flex;
      padding: 0 20px;
      align-items: center;
      color: #222;
      justify-content: space-between;
      .sum {
      }
    }
  }
  .playlist {
    width: 100%;
    height: auto;
    background-color: #fff;
  }
}
</style>