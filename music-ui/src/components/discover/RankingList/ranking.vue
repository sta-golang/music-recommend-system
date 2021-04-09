<template>
  <div>
    <!-- 榜单格式 -->
    <div class="ranking">
      <div class="title" :style="{background:'linear-gradient(to right,'+color[0]+','+color[1]+')'}">
        <div class="title_0">
          {{titleArr[0]}}
        </div>
        <div class="set">
          <div class="title_1">
            {{titleArr[1]}}
            {{titleArr[2]}}
          </div>
          <div class="timer">{{ time }}更新</div>
          <i class="iconfont">&#xe609;</i>
        </div>
      </div>
      <table class="centen" cellpadding='0'>
        <tr @click="rowClick(index,playlistDetail)" v-for="(item,index) in playlistDetail" :key='item.id'>
          <td :style="index < 3 ? 'color:red' : ''">{{index + 1}}</td>
          <td>{{item.name}}</td>
          <td>{{item.song.trim()}}</td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
// 格式化时间
import { formatDate } from '@/common/js/tool.js'
// 请求路由
import { _getSongsDetail, songDetail, _getMusicListDetail } from '@/network/discover/discover.js'
export default {
  props: ["rankId", 'title', 'color'],
  data () {
    return {
      playlistDetail: [],
      timer: 0
    }
  },
  created () {
    this.getPlaylistDetail()
  },
  methods: {
    // 获取歌单列表数据
    getPlaylistDetail () {
      _getMusicListDetail(this.rankId).then((result) => {
        if (result.code !== 200) return this.$message.error(result.msg)
        for (let i of result.playlist.tracks.splice(0, 8)) {
          _getSongsDetail(i.id).then(res => {
            let song = new songDetail(res.songs)
            this.playlistDetail.push(song);
          });
        }
        // console.log(this.playlistDetail);
        this.timer = result.playlist.updateTime
      })
    },
    // 点击了某一行
    rowClick (index, musiclist) {
      this.$bus.$emit('playMusic', index, musiclist)
    }
  },
  computed: {
    titleArr: function () {
      let titleclone = this.title
      return titleclone.substr(3).split('')
    },
    time () {

      return formatDate(new Date(this.timer), "MM月dd日");
    }
  }
}
</script>

<style lang='less' scoped>
.ranking {
  width: 334px;
  .title {
    position: relative;
    box-sizing: border-box;
    padding: 15px;
    display: flex;
    font-style: oblique;
    height: 90px;
    font-weight: 700;
    align-items: top;
    color: #fff;
    font-size: 20px;
    .title_0 {
      margin-top: -13px;
      font-size: 50px;
      margin-right: 10px;
    }
    .set {
      .timer {
        font-size: 12px;
        color: #ccc;
        font-style: normal;
        font-weight: 400;
      }
      .iconfont {
        position: absolute;
        right: 10px;
        top: 50%;
        transform: translateY(-50%);
        font-size: 50px;
        font-weight: 300;
      }
    }
  }
  .centen {
    width: 100%;
    height: 100%;
    border: 0;
    border-spacing: 0;
    border: 1px solid #ddd;
    box-sizing: border-box;
    border-top: 0;
    tr {
      width: 100%;
      // height: 30px;
      td {
        color: #222;
      }
      td:nth-child(1) {
        padding: 10px;
      }
      td:nth-child(3) {
        color: #666;
        text-align: right;
        padding-right: 10px;
      }
    }

    & :hover {
      background-color: #ccc;
    }
  }
}
.centen tr:nth-child(2n) {
  background-color: #f5f5f7;
}
</style>
