<template>
  <div class="musiclistdetail">
    <!-- 这个是歌单的详情页面 -->
    <detail-base-info :baseinfolist="baseinfolist"></detail-base-info>
    <detail-btns @toggleList="toggleList" :list="list"></detail-btns>
    <table-list
      @tableScroll="tableScroll"
      @rowClick="rowClick"
      :tracklist="tracklist"
      v-if="activeName == 0"
    ></table-list>
    <comment-list
      @scrollLoad="scrollLoad"
      :id="id"
      :hotCommentList="hotCommentList"
      :commentlist="commentlist"
      v-else-if="activeName == 1"
    ></comment-list>
    <collector v-else :id="id"></collector>
  </div>
</template>

<script>
// 请求路由
import {
  _getMusicListDetail,
  _getHotCommentlist,
  baseInfo,
  _getSongsDetail,
  songDetail,
  _getCommentlist
} from "@/network/discover/discover";
// 歌单介绍
import DetailBaseInfo from "./childrenComps/DetailBaseInfo";
// 歌单切换按钮
import DetailBtns from "./childrenComps/DetailBtns";
// 歌单列表
import TableList from "@/components/centent/musiclist/playlistMusicList";
// 歌单评论
import CommentList from "./childrenComps/CommentList";
// 歌单收藏者
import Collector from "./childrenComps/Collector";
export default {
  data() {
    return {
      playlist: [],
      // 歌单详情
      baseinfolist: [],
      // 按钮列表
      list: [],
      // 选中按钮
      activeName: 0,
      // 歌曲列表
      tracklist: [],
      // 剩下的歌曲列表
      leftList: [],
      // 歌单id
      id: "",
      limit: 10,
      offset: 1,
      // 评论列表
      commentlist: [],
      flag: true,
      // 最热评论
      hotCommentList: []
    };
  },
  methods: {
    // 子组件切换了
    toggleList(i) {
      this.activeName = i;
    },
    // 评论到底
    scrollLoad() {
      this.getCommentlist();
    },
    getCommentlist() {
      // 获取最新评论信息
      if (!this.flag) return;
      _getCommentlist({
        id: this.id,
        offset: this.offset,
        limit: this.limit
      }).then(result => {
        // console.log(result);
        if (result.comments.length == 0) {
          return (this.flag = false);
        }
        this.commentlist.push(...result.comments);
        this.offset += this.limit;
      });
      if (!this.hotCommentList.length) {
        this.getHotCommentlist();
      }
    },
    getHotCommentlist() {
      // 获取最热评论信息
      _getHotCommentlist({ id: this.id, type: 2, limit: 15 }).then(result => {
        // console.log(result);
        this.hotCommentList.push(...result.hotComments);
      });
    },
    getMusicListDetail() {
      // 如果这个歌单很多 会很卡 接口没有提供分页的功能
      // todos: 想到一个思路  是如果大于两百条就截取 只请求前两百条数据
      // 如果表格到底 v-infinite-scroll  通过这个事件来告诉父组件
      // 吧我们刚刚截取的数组剩下的去发请求
      console.log("getMusic-new");
      _getMusicListDetail(this.$route.query.id).then(result => {
        if (result.code !== 0) {
          return this.$message.error(result.message);
        }
        this.baseinfolist = new baseInfo(result.data);
        //let str = "评论(" + result.playlist.commentCount + ")";
        //this.list = ["歌曲列表", str, "收藏者"];
        //this.trackIds = result.playlist.trackIds;
        //this.playlist = result.playlist;
        console.log(result.data.user);
        this.tracklist = result.data.musics;
        console.log(this.tracklist);
        this.id = result.data.playlist.id;
        // 获取歌曲列表信息
        // 获取评论数据
        // this.getCommentlist();
      });
    },
    rowClick(index, list) {
      this.$bus.$emit("playMusic", index, list);
    },
    tableScroll() {
      if (this.leftList === []) return this.$message.info("没有更多了");
      for (let i of this.leftList) {
        _getSongsDetail(i.id).then(res => {
          let song = new songDetail(res.songs);
          this.tracklist.push(song);
        });
      }
      this.leftList = [];
    }
  },
  mounted() {
    // 获取歌单的信息
    this.getMusicListDetail();
  },
  components: {
    DetailBaseInfo,
    DetailBtns,
    TableList,
    CommentList,
    Collector
  },
  watch: {
    $route(to, from) {
      if (from.path === "/home/musiclistdetail") {
        this.tracklist = [];
        // console.log(to);
        // console.log(from);
        this.getMusicListDetail();
      }
    }
  }
};
</script>

<style lang="less" scpoped>
.musiclistdetail {
  margin-bottom: 50px;
}
</style>
