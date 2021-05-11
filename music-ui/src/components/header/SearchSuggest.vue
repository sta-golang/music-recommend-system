<template>
  <div class="searchSuggest">
    <div class="searchtext" @click="search()">
      搜“<span>{{ content }}</span
      >”相关的结果 >
    </div>
    <dl v-if="suggestList.songs">
      <dt><i class="iconfont">&#xe63f;</i>单曲</dt>
      <dd
        @click="songsClickHandle(item)"
        v-for="item in suggestList.songs"
        :key="item.id"
      >
        <span class="name">{{ item.name | sliceSplit }}</span>
        <span
          class="alias"
          v-for="(item2, index2) in item.alias"
          :key="index2"
          >{{ item2 }}</span
        >
        <span v-for="nameItem in item.artists" :key="nameItem.name"
          >- {{ nameItem.name }}</span
        >
      </dd>
    </dl>
    <dl v-if="suggestList.artists">
      <dt><i class="iconfont">&#xe604;</i>歌手</dt>
      <dd
        @click="artistsClickHandle(item)"
        v-for="item in suggestList.creators"
        :key="item.id"
      >
        <span class="name">{{ item.name }}</span>
        <span
          class="alias"
          v-for="(item2, index2) in item.alias"
          :key="index2"
          >{{ item2 }}</span
        >
        <span v-for="nameItem in item.artists" :key="nameItem.name"
          >- {{ nameItem.name }}</span
        >
      </dd>
    </dl>
  </div>
</template>

<script>
import { _getSongsDetail, songDetail } from "@/network/discover/discover";
export default {
  // 只有点击了单曲才是播放
  // 其他都是跳转到对应的界面
  props: {
    suggestList: {
      type: Object
    },
    content: {
      type: String
    }
  },
  filters: {
    sliceSplit(value) {
      if (!value) {
        return "";
      }
      if (value.length > 30) {
        return value.slice(0, 31) + "...";
      }
      return value;
    }
  },
  methods: {
    // 点击单曲
    songsClickHandle(item) {
      // 需要获取歌曲详情
      console.log(item);
      _getSongsDetail(item.id).then(res => {
        let song = new songDetail(res.data);
        console.log(song);
        console.log("song");
        this.$bus.$emit("pushPlayMusic", song);
      });
    },
    // 点击歌手
    artistsClickHandle(item) {
      this.$router.push({
        path: "/home/artistalbum",
        query: {
          id: item.id
        }
      });
    },
    playlistsClickHandle(item) {
      this.$router.push({
        path: "/home/musiclistdetail",
        query: {
          id: item.id
        }
      });
    },
    // 点击搜索
    search() {
      this.$router.push({
        path: "/home/searchlist",
        query: {
          content: this.content
        }
      });
    }
  },
  created() {
    console.log(this.suggestList);
  }
};
</script>

<style lang="less" scoped>
.searchSuggest {
  width: 300px;
  background-color: #f5f5f7;
  position: absolute;
  left: 0;
  top: 41px;
  z-index: 110;
  box-shadow: 4px 4px 20px rgba(0, 0, 0, 0.3);
  font-size: 10 px;
  .searchtext {
    cursor: pointer;
    height: 28px;
    line-height: 28px;
    padding-left: 10px;
    span {
      color: #0c73c2;
    }
  }
  dl {
    dt {
      color: #222;
      height: 28;
      line-height: 28px;
      background-color: #f5f5f7;
      i {
        display: inline-block;
        width: 30px;
        box-sizing: border-box;
        padding-right: 5px;
        text-align: right;
        line-height: 28px;
      }
    }
    dd {
      padding-left: 30px;
      height: 28px;
      line-height: 28px;
      background-color: #fafafa;
    }
    dd:hover {
      background-color: #ededed;
    }
    .name {
      color: #0c73c2;
    }
  }
}
</style>
