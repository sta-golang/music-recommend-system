<template>
  <div class="NewSong">
    <ul>
      <li v-for="(item,index) in newsonglist" @click="rowClick(index)" :key="item.id">
        <div class="number">{{index +1}}</div>
        <img :src="item.pic" alt="">
        <div class="desc">
          <p>{{item.name}}</p>
          <div class="author">
            <p>{{item.song }}</p>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { _getPersonalizedNewsong, _getSongsDetail, songDetail } from '@/network/discover/discover'
export default {
  data () {
    return {
      newsonglist: []
    }
  },
  methods: {
    getNewsongList () {
      _getPersonalizedNewsong().then(res => {
        for (let i of res.result) {
          _getSongsDetail(i.id).then(result => {
            let song = new songDetail(result.songs)
            this.newsonglist.push(song)
          })
        }
      })
    },
    // 点击了li
    rowClick (index) {
      // console.log(123);
      this.$bus.$emit('playMusic', index, this.newsonglist)
    },

  },
  created () {
    this.getNewsongList()
  }
}
</script>

<style lang='less' scoped>
.NewSong {
  ul {
    margin-bottom: 110px;
    display: flex;
    flex-wrap: wrap;
    border-left: 1px solid #eee;
    border-bottom: 1px solid #eee;
    li {
      cursor: pointer;
      border-top: 1px solid #eee;
      border-right: 1px solid #eee;
      box-sizing: border-box;
      padding: 10px;
      flex: 50%;
      display: flex;
      img {
        display: block;
        width: 42px;
        height: 42px;
      }
      .number {
        width: 30px;
        height: 42px;
        line-height: 42px;
        text-align: center;
      }
      .desc {
        margin-left: 10px;
        display: flex;
        flex-direction: column;
      }
    }
    & :hover {
      background-color: #ccc !important;
    }
    li:nth-child(4n + 1) {
      background-color: #f6f6f6;
    }
    li:nth-child(4n + 2) {
      background-color: #f6f6f6;
    }
  }
}
</style>
