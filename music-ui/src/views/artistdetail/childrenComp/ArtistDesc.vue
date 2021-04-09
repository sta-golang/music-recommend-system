<template>
  <div class="artistDesc">
    <div class="title">{{ artist.name }}</div>
    <!-- <p v-html="item.txt" style="white-space: pre-line;"></p> -->
    <p style="white-space: pre-line;">{{ artist.description }}</p>
  </div>
</template>

<script>
export default {
  props: {
    artist: {
      type: Object
    }
  },
  data() {
    return {};
  },
  methods: {},
  mounted() {
    // console.log(this.artist);
  }
};
</script>

<style lang="less" scoped>
.artistDesc {
  padding: 0 20px;
  .title {
    margin-top: 40px;
    margin-bottom: 10px;
    font-size: 16px;
    font-weight: 600;
    color: #222;
  }
  p {
    // text-indent: 2em;
    line-height: 30px;
    font-size: 14px;
    color: #666666;
  }
}
</style>
