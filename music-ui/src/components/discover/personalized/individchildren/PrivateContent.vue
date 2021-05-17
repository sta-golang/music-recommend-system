<template>
  <div class="PrivateContent">
    <ul>
      <li v-for="item in exclusiveBroadcastList" :key="item.id">
        <i class="iconfont">&#xe614;</i>
        <img :src="item.sPicUrl" alt="" />
        <p>{{ item.name }}</p>
      </li>
    </ul>
  </div>
</template>

<script>
import { _getExclusiveBroadcastList } from "@/network/discover/discover";
export default {
  data() {
    return {
      // 独家放送
      exclusiveBroadcastList: []
    };
  },
  methods: {
    // 获取独家放送列表
    async getExclusiveBroadcastList() {}
  },
  created() {
    this.getExclusiveBroadcastList();
  }
};
</script>

<style lang="less" scoped>
.PrivateContent {
  overflow: hidden;
  margin-bottom: 40px;
  ul {
    display: flex;
    justify-content: space-between;
    li {
      cursor: pointer;
      position: relative;
      width: 335px;
      height: auto;
      cursor: pointer;
      color: #222;
      i {
        position: absolute;
        color: #fff;
        left: 5px;
        font-size: 24px;
      }
      img {
        width: 100%;
      }
      p {
        margin-top: 2px;
      }
    }
  }
}
</style>
