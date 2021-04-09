<template>
  <div class="commentList" v-infinite-scroll='scrollLoad'>
    <div class="commentBox">
      <div class="area">
        <textarea v-model="content" @keyup.enter="addCommont" name="" id="" cols="30" rows="10"></textarea>
      </div>
      <span class="sub" @click="addCommont">评论</span>
    </div>
    <div class="title" v-if=" hotCommentList.length">
      最热评论
    </div>
    <ul v-if=" hotCommentList.length">
      <li class="item" v-for="(item,index) in hotCommentList" :key="index">
        <div class="image">
          <img :src="item.user.avatarUrl" alt="">
        </div>
        <div class="author">
          <span class="nickname">{{item.user.nickname}}:</span>
          <span class="content">{{item.content}}</span>
          <div class="beReplied" v-if="item.beReplied.length">
            <ul>
              <ol v-for="(item2,index2) in item.beReplied" :key="index2">
                <span class="user">
                  @{{item2.user.nickname}}
                </span>
                <span class="contont">
                  {{item2.content}}
                </span>
              </ol>
            </ul>
          </div>
          <div class="time">{{item.time | timer}}
            <div class="warp">
              <span>分享</span>
              <span>回复</span>
            </div>
          </div>
        </div>

      </li>
    </ul>
    <div class="title" v-if=" commentlist.length">
      最新评论
    </div>
    <ul>
      <li class="item" v-for="(item,index) in commentlist" :key="index">
        <div class="image">
          <img :src="item.user.avatarUrl" alt="">
        </div>
        <div class="author">
          <span class="nickname">{{item.user.nickname}}:</span>
          <span class="content">{{item.content}}</span>
          <div v-if=" item.beReplied   ">
            <div class="beReplied" v-if="item.beReplied.length">
              <ul>
                <ol v-for="(item2,index2) in item.beReplied" :key="index2">
                  <span class="user">
                    @{{item2.user.nickname}}
                  </span>
                  <span class="contont">
                    {{item2.content}}
                  </span>
                </ol>
              </ul>
            </div>
          </div>
          <div class="time">{{item.time | timer}}
            <div class="warp">
              <span>分享</span>
              <span>回复</span>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { formatDate } from '@/common/js/tool'
import { _SendComments } from '@/network/song'

export default {
  props: {
    commentlist: {
      type: Array,
      default: function () {
        return []
      }
    },
    id: {
      type: Number,
      default: function () {
        return 0
      }
    },
    hotCommentList: {
      type: Array,
      default: function () {
        return []
      }
    },
    commentType: {
      type: Number,
      default: function () {
        return 2
      }
    }
  },
  data () {
    return {
      content: ''
    }
  },
  methods: {
    scrollLoad () {
      // console.log('到底了');
      this.$emit('scrollLoad')
    },
    addCommont () {
      let cookie = this.$store.state.cookie
      if (cookie == '' || cookie == null) {
        return this.$message.error('请先登录才能评论！')
      }
      if (this.content.trim().length === 0) {
        this.content = this.content.trim()
        return this.$message.error('请输入评论内容')
      }
      if (this.content.trim().length >= 140) {
        return this.$message.error('不能超过140字')
      }
      _SendComments({
        type: this.commentType,
        t: 1,
        id: this.id,
        content: this.content,
        cookie
      }).then(result => {
        this.$message.success('评论成功')
        this.content = ''
        this.commentlist.unshift(result.comment)
      }).catch(err => {
        return this.$message.error('评论失败')
      })
    }
  },
  filters: {
    timer (value) {
      return formatDate(new Date(value), "yy年MM月dd日 hh-mm")
    }
  }
}
</script>

<style lang='less' scoped>
.commentList {
  margin-top: 20px;
  padding: 0 20px;
  .commentBox {
    overflow: hidden;
    box-sizing: border-box;
    padding: 10px;
    padding-bottom: 15px;
    background-color: #f0f0f2;
    textarea {
      resize: none;
      outline: none;
      padding: 5px;
      width: 98%;
      // width: 100% - 20px;
      height: 50px;
      max-height: 60px;
    }
    .sub {
      margin-top: 3px;
      // margin-bottom: 10px;
      border: 1px solid #ccc;
      padding: 3px 10px;
      float: right;
      background-color: #fff;
    }
  }
  .title {
    overflow: hidden;
    margin-top: 40px;
    color: #222;
    padding-bottom: 5px;
  }
  .item {
    border-top: 1px solid #ccc;
    padding: 17px 0;
    align-content: center;
    display: flex;
    .image {
      img {
        width: 40px;
        height: 40px;
        border-radius: 50%;
      }
    }

    .author {
      width: 100%;
      margin-left: 20px;
      .nickname {
        color: #0c73c2;
      }

      .content {
      }
      .beReplied {
        margin: 5px 0;
        padding: 5px 0;
        border-radius: 3px;
        background-color: #f1f1f4;
      }
      .time {
        width: 100%;
        position: relative;
        margin-top: 3px;
        .warp {
          position: absolute;
          bottom: 0;
          right: 10px;
          span {
            margin-left: 8px;
            padding-left: 5px;
            border-left: 1px solid #ccc;
          }
          span:first-child {
            border: 0 !important;
          }
        }
      }
    }
  }
}
</style>