<template>
  <div class="simiArtist">
    <div class="item" v-for="item in artist.similar_creator" :key="item.id">
      <img :src="item.image_url" alt="" @click="artistClick(item.id)" />
      <span>{{ item.name }}</span>
    </div>
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
  methods: {
    artistClick(id) {
      // console.log(id);
      this.$router.push({ path: "/home/artistalbum", query: { id } });
    }
  },
  mounted() {
    // console.log(this.artist);
  }
};
</script>

<style lang="less" scoped>
.simiArtist {
  .item {
    float: left;
    font-size: 18px;
    margin: 37px;
    margin-bottom: 3px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    img {
      width: 180px;
      height: 180px;
      border: 1px solid rgba(0, 0, 0, 0.3);
    }
    span {
      text-align: center;
      color: #333;
    }
  }
}
</style>
