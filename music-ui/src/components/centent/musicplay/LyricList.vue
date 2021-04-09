<template>
  <div class="lyriclist">
    <div>
      <!-- :style="{ transform: transformNum}" 滚动歌词这里想不到了 -->
      <ul v-if="lyricList.length">
        <li :class="currentItem==index?'currentItem':''" v-for="(item,index) in lyricList" :key="index">{{item}}</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    lyric: {
      type: String,
      default () {
        return ''
      }
    }
  },
  data () {
    return {
      lyricList: [],
      timerList: [],
      currentItem: 0,
    }
  },
  methods: {
    initLyric () {
      let str = '\n'
      let lyricList1 = this.lyric.split('\n').join('').split('[').slice(1)
      for (let i of lyricList1) {
        let tr = i.split(']')
        if (tr[1] == '') continue
        tr[0] = tr[0].split('.')[0];
        let timer1 = parseInt(tr[0].split(':')[0] * 60)
        let timer2 = parseInt(tr[0].split(':')[1])
        tr[0] = new Date(parseFloat(Number(timer1 + timer2)) * 1000)
        this.lyricList.push(tr[1])
        this.timerList.push(tr[0])
      }
      // console.log(this.lyricList)
      // console.log(this.timerList);
    }
  },
  mounted () {
    this.initLyric()
  },
  computed: {
    isMusicPlay () {
      return this.$store.state.isMusicPlay
    },
    currentTime () {
      return this.$store.state.currentTime
    },
    transformNum () {
      return 'translateY(' + this.currentItem * -20 + 'px)'
    }
  },
  watch: {
    lyric (val) {
      this.lyricList = []
      this.timerList = []
      this.initLyric()
    },
    isMusicPlay (val) {
      console.log(val);
    },
    currentTime (val) {
      this.timerList.forEach((value, index) => {
        //  后一项大于当前时间并且当前项小于当前时间 
        if (this.timerList[index + 1] >= val && val >= value) {
          this.currentItem = index
        }
      });
    }
  }
}
</script>

<style scoped lang='less'>
.lyriclist {
  margin-top: 10px;
  overflow: auto;
  height: 340px;
  ul {
    transition: 0.5s all;
    transform: translateY(0px);
    li {
      height: 28px;
      line-height: 28px;
      color: #333;
    }
    li.currentItem {
      font-size: 13px;
      font-weight: 600;
      color: #c62f2f;
    }
  }
}
</style>