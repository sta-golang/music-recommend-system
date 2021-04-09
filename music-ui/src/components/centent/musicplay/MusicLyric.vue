<template>
  <div class="musiclyric">
    <el-drawer ref='drawer' :modal='false' direction="btt" :visible.sync="showLyric" :with-header="false">
      <div class="drawerBox">
        <div class="lyricbox">
          <div class="img_left">
            <img :src="songDetailObj.pic" alt="">
            <div class="wrap">
              <span><i class="iconfont">&#xe60a;</i>喜欢</span>
              <span><i class="iconfont">&#xe63a;</i>收藏</span>
              <span><i class="iconfont">&#xe723;</i>VIP 下载</span>
              <span><i class="iconfont">&#xe65d;</i>分享</span>
            </div>
          </div>
          <div class="lyric_list">
            <i class="iconfont narrow" @click="closeLyricbox">&#xf01b4;</i>
            <div class="title">
              {{songDetailObj.name}}
            </div>
            <div class="songInfo">
              <div class="sognAlbum">
                专辑：<span>{{songDetailObj.album}}</span>
              </div>
              <div class="sognAlbum">
                歌手：<span>{{songDetailObj.song}}</span>
              </div>
            </div>
            <lyric-list :lyric='lyric'></lyric-list>
          </div>
        </div>
        <div class="commtentList">
          <comment-list @scrollLoad='scrollLoad' :id='songDetailObj.id' :hotCommentList='hotCommentList' :commentlist='commentList' :commentType='0'></comment-list>
        </div>
      </div>
      <div class="bgimage">
        <img :src="songDetailObj.pic" alt="">
      </div>
    </el-drawer>
  </div>
</template>

<script>
import { _getLyric, _getSongComment } from '@/network/song'
import CommentList from '@/views/musiclistdetail/childrenComps/CommentList'
import LyricList from './LyricList'
export default {
  data () {
    return {
      songDetailObj: {},
      // 歌词
      lyric: '',
      commentOffset: 0,
      hotCommentList: [],
      commentList: [],
      commentlimit: 20
    }
  },
  methods: {
    closeLyricbox () {
      this.$store.commit('editshowLyric', false)
    },
    getLyric (id) {
      _getLyric(id).then(result => {
        // console.log(result);
        this.lyric = result.lrc.lyric ? result.lrc.lyric : '纯音乐'
        // console.log(this.lyric);
      })
    },
    getSongComment () {
      _getSongComment({
        id: this.songDetailObj.id,
        offset: this.commentOffset,
        limit: this.commentlimit
      }).then(result => {
        console.log(result);
        if (this.hotCommentList.length === 0) {
          this.hotCommentList = result.hotComments
        }
        this.commentList.push(...result.comments)
        this.commentOffset += this.commentlimit
      })
    },
    scrollLoad () {
      this.getSongComment()
    }
  },
  computed: {
    showLyric () {
      return this.$store.state.showLyric
    },
    songDetail () {
      return this.$store.state.songDetail
    }
  },
  watch: {
    songDetail: function (obj) {
      this.songDetailObj = obj
      this.getLyric(obj.id)
      this.getSongComment()
      // console.log(obj);
    }
  },
  components: {
    LyricList,
    CommentList
  }
}
</script>

<style lang='less'>
i {
  cursor: pointer;
}
.el-drawer__open .el-drawer.btt {
  height: calc(100% - 50px) !important;
}
.musiclyric {
  .el-drawer {
    height: calc(100% - 50px - 60px);
    overflow: scroll;
    .drawerBox {
      position: absolute;
      width: 866px;
      top: 0;
      left: 0;
      right: 0;
      margin: auto;
      z-index: 260;
      .lyricbox {
        display: flex;
        justify-content: space-between;
        .img_left {
          width: 300px;
          height: 500px;
          display: flex;
          flex-direction: column;
          align-items: center;
          img {
            margin-top: 70px;
            width: 200px;
            height: 200px;
            border-radius: 50%;
            border: 50px solid #222;
            box-shadow: 0px 0px 10px 10px rgba(255, 255, 255, 0.7);
          }
          .wrap {
            margin-top: 30px;
            color: #333;
            width: 100%;
            display: flex;
            flex-direction: row;
            justify-content: space-between;
            span {
              height: 26px;
              line-height: 26px;
              padding: 0 10px;
              background-color: #f5f5f6;
              border: 1px solid #acadae;
              border-radius: 5px;
              i {
                font-size: 14px;
                margin-right: 5px;
              }
            }
          }
        }
        .lyric_list {
          padding-top: 20px;
          box-sizing: border-box;
          position: relative;
          margin-left: 20px;
          width: 55%;
          height: 500px;
          .narrow {
            position: absolute;
            top: 20px;
            width: 35px;
            right: 0;
            height: 25px;
            color: #666;
            text-align: center;
            font-size: 20px;
            border: 1px solid #d1d1d2;
            background-color: rgba(255, 255, 255, 0.9);
          }
          .narrow:hover {
            color: #333;
          }
          .title {
            color: #222;
            font-size: 24px;
            margin: 20px 0;
          }
          .songInfo {
            display: flex;
            font-size: 13px;
            color: #222;
            div {
              margin-right: 20px;
              span {
                color: #1350ac;
              }
              span:hover {
                color: #0e3f88;
              }
            }
          }
        }
      }
      .commtentList {
        margin-left: -30px;
        margin-bottom: 50px;
        width: 60%;
      }
    }
    .bgimage {
      height: 100%;
      width: 100%;
      text-align: center;
      img {
        -webkit-filter: blur(100px);
        filter: blur(150px);
        margin-top: -50px;
        width: 50%;
        height: 400px;
      }
    }
  }
}
</style>