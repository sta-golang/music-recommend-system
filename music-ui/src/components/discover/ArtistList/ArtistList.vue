<template>
  <div
    class="ArtistList"
    v-infinite-scroll="scrollLoad"
    infinite-scroll-immediate
    infinite-scroll-delay="1000"
  >
    <div class="title">
      <div class="area">
        语种：
        <span
          v-for="(item, index) in area"
          @click="areaClick(index)"
          :key="item.value"
          :class="areaIndex === index ? 'current' : ''"
          >{{ item.name }}</span
        >
      </div>
      <div class="type">
        分类：
        <span
          v-for="(item, index) in type"
          @click="typeClick(index)"
          :key="item.value"
          :class="typeIndex === index ? 'current' : ''"
          >{{ item.name }}</span
        >
      </div>
    </div>
    <div class="scroll">
      <artist :artistlist="artistlist"></artist>
    </div>
  </div>
</template>

<script>
// 请求路由
import { _getArtistlist } from "@/network/discover/discover";
// import Musiclist from '@/components/centent/musiclist/Musiclist'
import Artist from "./Artist";
export default {
  data() {
    return {
      // 当前地区
      areaIndex: 0,
      // 当前分类
      typeIndex: 0,
      limit: 30,
      page: 1,
      offset: 0,
      artistlist: [],
      area: [
        { value: -1, name: "全部" },
        { value: 7, name: "华语" },
        { value: 96, name: "欧美" },
        { value: 8, name: "日本" },
        { value: 16, name: "韩国" },
        { value: 0, name: "其他" }
      ],
      type: [
        { value: -1, name: "全部" },
        { value: 1, name: "男歌手" },
        { value: 2, name: "女歌手" },
        { value: 3, name: "乐队" }
      ]
    };
  },
  methods: {
    scrollLoad() {
      // console.log('到底');
      this.getArtistlist();
    },
    areaClick(index) {
      this.offset = 0;
      this.artistlist = [];
      this.areaIndex = index;
      this.getArtistlist();
    },
    typeClick(index) {
      this.typeIndex = index;
      this.offset = 0;
      this.artistlist = [];
      this.getArtistlist();
    },
    getArtistlist() {
      _getArtistlist({
        page: this.page
      }).then(result => {
        // console.log(result);
        console.log(result.data);
        this.artistlist.push(...result.data);
        this.page += 1;
      });
    }
  },
  created() {},
  components: {
    Artist
  }
};
</script>

<style lang="less" scoped>
.ArtistList {
  .title {
    margin-top: 10px;
    font-size: 16px;
    border-bottom: 1px solid #ccc;
    padding-bottom: 5px;
    margin-bottom: 10px;
    font-size: 12px;
    overflow: hidden;
    span {
      margin: 0 10px;
      padding: 2px 5px;
      border-radius: 5px;
      cursor: pointer;
    }
    .current {
      color: #fff;
      background-color: #777;
    }
    .area {
      overflow: hidden;
      margin: 0 0 10px;
    }
    div {
      margin: 0 0 25px;
    }
  }
  .scroll {
  }
}
</style>
