<template>
  <div>
    <artist-base-info :artist="artist"></artist-base-info>
    <artist-btns :btns="btns" @artisttoggle="artisttoggle"></artist-btns>
    <artist-song v-if="itemIndex == 0" :artist="artist"></artist-song>
    <simi-artist v-else-if="itemIndex == 1" :artist="artist"></simi-artist>
  </div>
</template>

<script>
import {
  _getArtistalbumDesc,
  _getArtistalbum,
  _getArtistTop50
} from "@/network/discover/discover";
import ArtistBaseInfo from "./childrenComp/ArtistBaseInfo";
import ArtistBtns from "./childrenComp/ArtistBtns";
import ArtistAlbum from "./childrenComp/ArtistAlbum";
import SimiArtist from "./childrenComp/SimiArtist";
import ArtistSong from "./childrenComp/ArtistSong";
export default {
  data() {
    return {
      id: "",
      // 头部信息
      artist: {},
      // 切换按钮
      btns: ["歌曲", "相似歌手"],
      itemIndex: 0
    };
  },
  methods: {
    // 获取专辑信息
    getArtistalbum() {
      _getArtistalbum(this.id).then(result => {
        this.artist = result.data;
      });
    },
    // 切换
    artisttoggle(i) {
      this.itemIndex = i;
      // console.log(i);
    }
  },
  created() {
    this.id = this.$route.query.id;
    this.getArtistalbum();
  },
  components: {
    ArtistBaseInfo,
    ArtistAlbum,
    ArtistSong,
    SimiArtist,
    ArtistBtns
  },
  watch: {
    $route(to, from) {
      if (from.path === "/home/artistalbum") {
        this.id = this.$route.query.id;
        this.artist = {};
        this.getArtistalbum();
      }
    }
  }
};
</script>

<style></style>
