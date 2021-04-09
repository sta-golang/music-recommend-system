<template>
  <div class="Artist">
    <div class="musiclist">
      <ul>
        <li
          @click="artistClick(item.id)"
          v-for="item in artistlist"
          :key="item.id"
        >
          <img :src="item.image_url ? item.image_url : image_url" alt="" />
          <p>{{ item.name }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: ["artistlist"],
  data() {
    return {};
  },
  methods: {
    // 跳转到歌手详情页
    artistClick(id) {
      // console.log(id);
      this.$router.push({ path: "/home/artistalbum", query: { id } });
    }
  },
  created() {}
};
</script>

<style lang="less" scoped>
.Artist {
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
}
</style>
