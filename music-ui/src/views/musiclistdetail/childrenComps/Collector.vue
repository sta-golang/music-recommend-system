<template>
  <div v-infinite-scroll='scrollLoad' :infinite-scroll-delay='2000'>
    <ul class="collector" v-if="subscribersList.length">
      <li class="item" v-for="item in subscribersList" :key="item.userId">
        <img :src="item.avatarUrl" alt="">
        <span>{{item.nickname}}</span>
      </li>
    </ul>
  </div>
</template>

<script>
import { _getCollector } from '@/network/discover/discover'
export default {
  props: ['id'],
  data () {
    return {
      limit: 20,
      offset: 0,
      subscribersList: [],
      // 判断是否还有
      flag: true
    }
  },
  methods: {
    // 获取收藏者信息
    getCollector () {
      if (!this.flag) return
      _getCollector({
        id: this.id,
        limit: this.limit,
        offset: this.offset
      }).then(result => {
        if (result.subscribers.length === 0) return this.flag = false
        console.log(result);
        this.subscribersList.push(...result.subscribers)
      })
      this.offset += this.limit
    },
    scrollLoad () {
      // console.log('到底');
      this.getCollector()
    }
  },
  created () {
  }
}
</script>

<style lang='less' scoped>
.collector {
  display: flex;
  flex-wrap: wrap;
  flex-direction: row;
  .item {
    width: 130px;
    margin: 20px 40px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    img {
      width: 55px;
      height: 55px;
      border-radius: 50%;
      border: 1px solid rgba(0, 0, 0, 0.3);
    }
    span {
      margin-top: 5px;
      text-align: center;
    }
  }
}
</style>