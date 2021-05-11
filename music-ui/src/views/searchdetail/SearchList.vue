<template>
  <div class="searchList">
    <div class="searchTitle">
      搜索“<span>{{ content }}</span
      >”,找到<span>{{ sum }}</span
      >首单曲
    </div>
    <table-list @rowClick="rowClick" :tracklist="searchList"></table-list>
  </div>
</template>

<script>
import { _getSongsDetail, songDetail } from "@/network/discover/discover";
import { _getSearchList } from "@/network/search";
import TableList from "@/components/centent/musiclist/TableList";
export default {
  data() {
    return {
      content: "",
      searchList: []
    };
  },
  methods: {
    getSearchList(type) {
      _getSearchList(this.content, type).then(result => {
        for (let i of result.data) {
          console.log(i);
          this.searchList.push(i);
        }
      });
    },
    rowClick(index, list) {
      this.$bus.$emit("pushPlayMusic", list[index]);
    }
  },
  created() {
    this.content = this.$route.query.content;
    this.getSearchList(1);
  },
  computed: {
    sum() {
      return this.searchList.length;
    }
  },
  watch: {
    $route(to, from) {
      if (from.path === "/home/searchlist") {
        this.searchList = [];
        this.content = this.$route.query.content;
        this.getSearchList(1);
      }
    }
  },
  components: {
    TableList
  }
};
</script>

<style lang="less" scoped>
.searchTitle {
  margin: 20px;
  span {
    color: #1b8ee6;
  }
}
</style>

