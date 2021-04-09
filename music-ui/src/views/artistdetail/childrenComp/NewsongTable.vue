<template>
  <div class="newsongtable">
    <table cellpadding="0">
      <tr
        @click="rowClick(index, musiclist)"
        v-for="(item, index) in musiclist"
        :key="index"
      >
        <td>{{ index + 1 }}</td>
        <td><img :src="item.image_url" alt="" /></td>
        <td>{{ item.name }}</td>
        <td>{{ item.title }}</td>
        <td>{{ item.creator_names }}</td>
        <td>{{ item.play_time | formatDate }}</td>
      </tr>
    </table>
  </div>
</template>

<script>
import { formatDate } from "../../../common/js/tool.js";
export default {
  props: ["musiclist"],
  filters: {
    formatDate(time) {
      var data = new Date(time);
      return formatDate(data, "mm:ss");
    }
  },
  methods: {
    // 点击了某一行
    rowClick(index, musiclist) {
      this.$bus.$emit("playMusic", index, musiclist);
    }
  }
};
</script>

<style lang="less" scoped>
.newsongtable {
  table {
    width: 100%;
    border: none;
    border-spacing: 0;
    border: 1px solid #eee;
    tr {
      width: 100%;
      display: flex;
      height: 60px;
      td {
        display: flex;
        align-items: center;
      }
      td:nth-child(1) {
        color: #ccc;
        width: 50px;
        box-sizing: border-box;
        padding: 20px;
      }

      td:nth-child(2) {
        width: 50px;
        img {
          width: 50px;
          height: 50px;
        }
      }
      td:nth-child(3) {
        padding-left: 30px;
        flex: 7;
      }
      td:nth-child(4) {
        flex: 3;
      }
      td:nth-child(5) {
        flex: 1;
      }
    }
    tr:nth-child(2n) {
      background-color: #f5f5f7;
    }
  }
}
</style>
