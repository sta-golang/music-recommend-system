<template>
  <div class="RankingList">
    <div class="title">
      官方榜
    </div>
    <div class="wrap ">
      <ranking v-if="rankList[0]" :rank-id=' rankList[0].id' :color='["#3c8dde","#638fef"]' :title='  rankList[0].name'></ranking>
      <ranking v-if="rankList[1]" :rank-id=' rankList[1].id' :color='["#1499b4","#44c4cd"]' :title='  rankList[1].name'></ranking>
      <ranking v-if="rankList[2]" :rank-id=' rankList[2].id' :color='["#d0456e","#ee658d"]' title="123原创榜"></ranking>
      <ranking v-if="rankList[3]" :rank-id=' rankList[3].id' :color='["#c25a48","#b9613e"]' :title='  rankList[3].name'></ranking>
      <ranking v-if="rankList[4]" :rank-id=' rankList[4].id' :color='["#8237c2","#a960d3"]' :title='  rankList[4].name'></ranking>
    </div>
    <div class="title">
      全球榜
    </div>
    <musiclist @songListDetails='songListDetails' :RecommendedSongList='rankList.slice(5)'></musiclist>

  </div>
</template>

<script>
// 子组件
import ranking from './ranking'
// 列表组件
import musiclist from '@/components/centent/musiclist/Musiclist'
// 请求路由
import { _getTopListDetail } from '@/network/discover/discover'
export default {
  data () {
    return {
      rankList: []
    }
  },
  methods: {
    // 获取榜单列表
    async getToplist () {
      _getTopListDetail().then(result => {
        // console.log(result);
        if (result.name && result.name === 'Error') {
          return this.$message.error('请求错误')
        }
        if (result.code !== 200) return this.$message.error(result.msg)
        this.rankList = result.list
      })
    },
    songListDetails (id) {
      this.$router.push({ path: '/home/musiclistdetail', query: { id } })
    }
  },
  created () {
    this.getToplist()
  },
  components: {
    ranking,
    musiclist
  }
}
</script>

<style lang='less' scoped>
.RankingList {
  .title {
    margin-top: 20px;
    font-size: 16px;
    border-bottom: 1px solid #ccc;
    padding-bottom: 5px;
    margin-bottom: 20px;
  }
  .wrap {
    display: flex;
    flex-wrap: wrap;
    width: 1100px;
    div {
      margin-bottom: 20px;
      margin-right: 22px;
    }
  }
}
</style>
