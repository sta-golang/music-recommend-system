<template>
  <div class="ArtistBtns">
    <span :class="['artistItem',itemIndex ==index ? 'current' :'']" @click="artisttoggle(index)" v-for="(item,index) in btns" :key="index">{{item}}</span>
  </div>
</template>

<script>
export default {
  props: ["btns"],
  data () {
    return {
      itemIndex: 0
    }
  },
  methods: {
    artisttoggle (index) {
      this.itemIndex = index
      this.$emit('artisttoggle', index)
    }
  }
}
</script>

<style lang='less' scoped>
.ArtistBtns {
  width: 110%;
  margin-left: -20px;
  margin-top: 30px;
  border-bottom: 1px solid #ccc;
  padding-bottom: 5px;
  padding-left: 30px;
  .artistItem {
    cursor: pointer;
    font-size: 14px;
    margin: 30px;
    padding-bottom: 3px;
  }
  .current {
    color: #c62f2f;
    border-bottom: 3px solid #c62f2f;
  }
}
</style>