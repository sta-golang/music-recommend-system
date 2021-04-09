<template>
  <div class="DetailBtns">
    <div>
      <!-- 切换歌单详情 -->
      <span :class="activeName==index?'current':''" @click="spanClick(index)" v-for="(item,index) in list" :key="index">
        {{item}}
      </span>
    </div>
  </div>
</template>

<script>
export default {
  props: ['list'],
  data () {
    return {
      activeName: 0
    }
  },
  methods: {
    // 切换样式并通知父组件
    spanClick (i) {
      this.activeName = i
      this.$emit('toggleList', i)
    }
  }
}
</script>

<style lang='less' scoped>
.DetailBtns {
  width: 100%;
  margin-left: -15px;
  margin-top: 30px;
  border-bottom: 1px solid #ccc;
  padding-bottom: 5px;
  padding-left: 30px;
  div {
    padding: 0 0 5px;
    span {
      cursor: pointer;
      text-align: center;
      margin: 0 10px;
      padding: 0 2px 5px;
      font-size: 15px;
    }
    .current {
      border-bottom: 3px solid #c62f2f;
    }
  }
}
</style>