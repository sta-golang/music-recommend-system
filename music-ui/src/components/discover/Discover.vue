<template>
  <div>
    <div class="Discover">
      <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="个性推荐" name="first">
          <personalized v-if="activeName == 'first'"></personalized>
        </el-tab-pane>
        <el-tab-pane label="歌单" name="second">
          <popularsonglist v-if="activeName == 'second'"></popularsonglist>
        </el-tab-pane>
        <el-tab-pane label="主播电台" name="third">主播电台</el-tab-pane>
        <el-tab-pane label="排行榜" name="fourth">
          <ranking-list v-if="activeName == 'fourth'"></ranking-list>
        </el-tab-pane>
        <el-tab-pane label="歌手" name="fifth">
          <artist-list v-if="activeName == 'fifth'"></artist-list>
        </el-tab-pane>
        <el-tab-pane label="最新音乐" name="sixth">
          <newsong-list v-if="activeName == 'sixth'"></newsong-list>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
<script>
// 个性推荐
import Personalized from './personalized/Personalized'
// 热门歌单
import Popularsonglist from './Popularsonglist/Popularsonglist'
// 排行榜
import RankingList from './RankingList/RankingList'
// 歌手列表 
import ArtistList from './ArtistList/ArtistList'
// 最新音乐
import NewsongList from './newsongList/NewsongList'
export default {
  data () {
    return {
      activeName: 'first'
    }
  },
  methods: {
    handleClick () { },
  },
  components: {
    Popularsonglist,
    RankingList,
    ArtistList,
    NewsongList,
    Personalized
  },
  created () {
  }
}
</script>

<style lang='less' scoped>
.Discover {
  width: 100%;
  position: relative;
  margin: auto;
  left: 0;
  right: 0;
  bottom: 0;
  top: 0;
}
</style>
