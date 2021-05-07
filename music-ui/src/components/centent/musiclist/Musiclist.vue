<template>
  <div class="musiclist">
    <ul>
      <li
        v-for="item in RecommendedSongList"
        @click="songListDetails(item.id)"
        :key="item.id"
      >
        <img :src="item.image_url" alt="" />
        <p>{{ item.name }}</p>
        <div v-if="item.playCount" class="count">
          <i class="iconfont">&#xe6cf;</i>{{ item.playCount }}
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: ["RecommendedSongList"],
  data() {
    return {};
  },
  methods: {
    songListDetails(id) {
      this.$emit("songListDetails", id);
    }
  },
  created() {}
};
</script>

<style lang="less" scoped>
.musiclist {
  ul {
    display: flex;
    flex-wrap: wrap;
    margin-bottom: 20px;
    width: 1100px;
    li {
      cursor: pointer;
      margin-right: 17px;
      margin-bottom: 20px;
      position: relative;
      display: flex;
      flex-direction: column;
      width: 195px;
      color: #222;
      font-size: 14px;
      img {
        width: 195px;
        height: 195px;
      }
      .count {
        color: #fff;
        text-align: right;
        width: 195px;
        height: 22px;
        line-height: 22px;
        top: 0;
        right: 0;
        position: absolute;
        background-image: linear-gradient(
          to right,
          rgba(0, 0, 0, 0),
          rgba(0, 0, 0, 0.3)
        );

        padding-right: 5px;
        i {
          font-size: 12px;
          margin-right: 5px;
        }
      }
    }
  }
}
</style>
