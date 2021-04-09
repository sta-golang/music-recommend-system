<template>
  <div class="DetailBaseInfo">
    <!-- 歌单介绍 -->
    <div class="image">
      <img :src="baseinfolist.img" alt="">
    </div>
    <div class="info">
      <div class="title">{{baseinfolist.name}}</div>
      <div class="author">
        <span><img :src="baseinfolist.avatarUrl" alt=""></span>
        <span class="nickname">{{baseinfolist.nickname}}</span>
        <span class="timer">{{baseinfolist.createTime | timer}}创建</span>
      </div>
      <div class="btns">
        <div class="btns_1"><i class="iconfont">&#xe609; </i> 播放全部 <i class="iconfont">&#xe522;</i></div>
        <div class="btns_2"><i class="iconfont">&#xe63a; </i> 收藏( {{baseinfolist.subscribedCount}} )</div>
        <div class="btns_3"><i class="iconfont">&#xe65d; </i> 分享( {{baseinfolist.shareCount}} )</div>
        <div class="btns_3"><i class="iconfont">&#xe723; </i> 下载全部</div>
      </div>
      <div class="tags">
        标签：<span v-for="(item,index) in baseinfolist.tags" :key="index">{{item }}</span>
      </div>
      <div class="desc">
        <div>简介：</div> <span v-html="baseinfolist.desc" :class="control =='up'?'overhidden':''" style="white-space: pre-line;"></span>
        <div class="control iconfont" @click="controlHandle('down')" v-if="control == 'up'">&#xe65b;</div>
        <div class="control iconfont" @click="controlHandle('up')" v-else>&#xe659;</div>
      </div>
      <div class="numwrap">
        <span>歌曲数<p>{{baseinfolist.trackCount}}</p></span>
        <span>播放数<p>{{baseinfolist.playCount | bignum}}</p></span>
      </div>

    </div>
  </div>
</template>

<script>
import { formatDate, bignumSlice } from '@/common/js/tool.js'
export default {
  props: ["baseinfolist"],
  data () {
    return {
      control: 'up'
    }
  },
  methods: {
    // 切换简介展示
    controlHandle (val) {
      this.control = val
    }
  },
  mounted () {
    // console.log(this.baseinfolist);
  },
  filters: {
    timer (val) {
      return formatDate(new Date(val), "yy-MM-dd")
    },
    bignum (val) {
      return bignumSlice(val)
    }
  }
}
</script>

<style lang='less' scoped>
.DetailBaseInfo {
  margin-top: 10px;
  padding: 10px;
  display: flex;
  .image {
    img {
      width: 200px;
      height: 200px;
    }
  }

  .info {
    margin-left: 20px;
    flex-direction: column;
    position: relative;
    width: 100%;
    .title {
      font-size: 20px;
      color: #111;
    }
    .author {
      margin: 10px 0;
      display: flex;
      align-items: center;
      span {
        img {
          width: 30px;
          height: 30px;
          border-radius: 50%;
        }
        margin-right: 5px;
      }
      .nickname {
        font-size: 14px;
        margin-right: 20px;
      }
      .timer {
        color: #666;
        font-size: 12px;
      }
    }

    .btns {
      display: flex;
      margin: 20px 0 18px;
      font-size: 12px;
      div {
        color: #444;
        margin-right: 10px;
        border-radius: 4px;
        padding: 0px 11px;
        height: 28px;
        line-height: 28px;
        border: 1px solid #ccc;
      }
      .btns_1 {
        color: #fff;
        border: 1px solid #c62f2f;
        background-color: #c62f2f;
        i {
          margin-left: 5px;
          font-size: 12px;
        }
      }

      .btns_2 {
      }

      .btns_3 {
      }
    }

    .tags {
      margin-bottom: 5px;
      font-size: 12px;
      color: #222;
      span {
        margin-left: 5px;
        color: #0c73c2;
      }
    }

    .desc {
      color: #222;
      display: flex;
      width: 800px;
      div {
        overflow: hidden;
        width: 40px;
      }
      span {
        width: 800px;
        display: inline-block;
      }
      .overhidden {
        height: 40px;
        display: -webkit-box;
        display: -moz-box;
        white-space: pre-wrap;
        word-wrap: break-word;
        overflow: hidden;
        text-overflow: ellipsis;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 2;
      }
    }
    .numwrap {
      position: absolute;
      top: 0;
      right: 20px;
      span {
        display: inline-block;
        text-align: center;
        padding: 0 10px;
      }
      span:first-child {
        border-right: 1px solid #ccc;
      }
    }
    .control {
      position: absolute;
      top: 170px;
      right: 30px;
      color: #222;
      font-size: 26px;
      cursor: pointer;
    }
  }
}
</style>