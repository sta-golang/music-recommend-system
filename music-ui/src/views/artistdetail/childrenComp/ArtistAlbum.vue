<template>
  <div class="artistAlbum">
    <div class="item" v-for="item in albumList" :key="item.id">
      <div class="image">
        <img :src="item.picUrl" alt="">
        <div class="timer">{{item.publishTime | timer}}</div>
      </div>
      <div class="sheet">
        <div class="title">{{item.name}}
          <span>
            <i class="iconfont">&#xe63a;</i>
            <i class="iconfont">&#xe609;</i>
          </span>
        </div>
        <!-- {{item}} -->
        <artist-album-list class="table" :id="item.id"></artist-album-list>
      </div>
    </div>
  </div>
</template>

<script>
import { _getArtistalbum } from '@/network/discover/discover'
import ArtistAlbumList from './ArtistAlbumList'
import { formatDate } from '@/common/js/tool.js'
export default {
  props: ["id"],
  data () {
    return {
      // 专辑列表
      albumList: [123],
    }
  },
  methods: {
    getArtistalbum () {
      _getArtistalbum(this.id).then(result => {
        this.albumList = result.hotAlbums
        // console.log(this.albumList);
      })
    },
  },
  mounted () {
    this.getArtistalbum()
  },
  components: {
    ArtistAlbumList
  },
  filters: {
    timer (val) {
      return formatDate(new Date(val), "yy-MM-dd")
    }
  },
  watch: {
    '$route' (to, from) {
      if (from.path === '/home/artistalbum') {
        this.id = this.$route.query.id
        this.getArtistalbum()
      }
    }
  }
}
</script>

<style lang='less' scoped>
.artistAlbum {
  margin-top: 10px;
  overflow: hidden;
  .item {
    display: flex;
    margin-top: 10px;
    margin-bottom: 40px;
    .image {
      position: relative;
      width: 150px;
      img {
        display: inline-block;
        width: 150px;
        height: 150px;
        position: relative;
        z-index: 2;
      }
    }
    .image::after {
      content: "";
      height: 130px;
      width: 150px;
      position: absolute;
      top: 12px;
      left: 20px;
      z-index: 0;
      background: url("../../../common/img/coverall.png") no-repeat;
      background-position: 0 -848px;
    }
    // .sheet::after {
    //   content: "";
    //   display: inline-block;
    //   position: absolute;
    //   width: 150px;
    //   top: 0;
    //   z-index: 0;
    //   left: -190px;
    //   height: 150px;
    //   border-radius: 50%;
    //   background-color: #000;
    // }
    // .image::before {
    //   content: "";
    //   display: inline-block;
    //   position: absolute;
    //   width: 10px;
    //   top: 10px;
    //   z-index: 1;
    //   right: -10px;
    //   height: 50px;
    //   background-color: rgba(255, 255, 255, 0.2);
    //   border: 1px solid rgba(0, 0, 0, 0.2);
    // }
    // .image::after {
    //   content: "";
    //   display: inline-block;
    //   position: absolute;
    //   width: 10px;
    //   top: 80px;
    //   z-index: 1;
    //   right: -10px;
    //   height: 50px;
    //   background-color: rgba(255, 255, 255, 0.2);
    //   border: 1px solid rgba(0, 0, 0, 0.2);
    // }

    .sheet {
      position: relative;
      margin-left: 60px;
      width: 100%;
      .title {
        color: #222;
        margin-bottom: 10px;
        flex: 1;
        position: relative;
        span {
          position: absolute;
          right: 0;
          top: 0;
          height: 30px;
          display: flex;
          justify-content: center;
          i {
            margin-left: 20px;
          }
        }
      }
      .table {
        border: 1px solid #ccc;
      }
    }
  }
}
</style>