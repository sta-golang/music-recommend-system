<template>
  <div>
    <table-list @rowClick='rowClick' :tracklist='tracklist' :hiddenTableHeader='false' :hiddenSinger='false' :hiddenAlbum='false'></table-list>
  </div>
</template>

<script>
import TableList from '@/components/centent/musiclist/TableList'
import { _getAlbum, songDetail, _getSongsDetail } from '@/network/discover/discover'
export default {
  props: ["id"],
  data () {
    return {
      tracklist: []
    }
  },
  methods: {
    getAlbum () {
      // console.log(this.id);
      if (!this.id) return
      _getAlbum(this.id).then(result => {
        // console.log(result);
        for (let i of result.songs) {
          _getSongsDetail(i.id).then(res => {
            // console.log(res);
            this.tracklist.push(new songDetail(res.songs))
          })
        }
        // console.log(this.tracklist);
      })
    },
    rowClick (index, list) {
      this.$bus.$emit('playMusic', index, list)
    }
  },
  mounted () {
    this.getAlbum()
  },
  components: {
    TableList
  }
}
</script>

<style>
</style>