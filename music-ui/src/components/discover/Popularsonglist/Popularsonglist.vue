<template>
  <div class="Popularsonglist">
    <div>
      <musiclist
        @songListDetails="songListDetails"
        :RecommendedSongList="musicList"
      ></musiclist>
    </div>
  </div>
</template>

<script>
// 列表组件
import Musiclist from "@/components/centent/musiclist/Musiclist";
import {
  _getHotPlaylist,
  _getTopHighquality
} from "@/network/discover/discover";
export default {
  data() {
    return {
      tags: [],
      currentIndex: 0,
      limit: 24,
      page: 1,
      musicList: []
    };
  },
  methods: {
    getHotPlaylist() {
      _getHotPlaylist().then(result => {
        this.musicList = result.data;
      });
    },
    // 点击tag后修改并换list
    tagClick(index) {
      this.currentIndex = index;
      this.getMusicList();
    },
    // 获取对应标签的list
    getMusicList() {
      _getTopHighquality({
        cat: this.tags[this.currentIndex].name,
        limit: this.limit,
        time: new Date().getTime()
      }).then(result => {
        this.musicList = result.playlists;
      });
    },
    songListDetails(id) {
      this.$router.push({ path: "/home/musiclistdetail", query: { id } });
    }
  },
  created() {
    this.getHotPlaylist();
  },
  components: {
    Musiclist
  }
};
</script>

<style lang="less" scoped>
.Popularsonglist {
  .tags {
    display: flex;
    font-size: 14px;
    margin-bottom: 30px;
    span {
      color: #222;
      margin-right: 10px;
    }
    .tag-item {
      padding: 0 15px;
      color: #666;
      cursor: pointer;
      border-right: 1px solid #ccc;
    }
    & :last-child {
      border: none;
    }
    .current {
      color: #222;
      font-weight: 600;
    }
  }
}
</style>
