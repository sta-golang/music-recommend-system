<template>
  <div>
    <div class="musicplay">
      <div class="play-top" v-show="showTop">
        <div class="image">
          <img :src="userinfo.pic" alt="" />
          <i class="iconfont imageIconfont" @click="showLyricBox">&#xe625;</i>
        </div>
        <div class="userinfo">
          <div class="info">
            <div class="title">{{ userinfo.name }}</div>
            <div class="nickname">{{ userinfo.song }}</div>
          </div>
          <div class="warp">
            <i class="iconfont">&#xe60a;</i>
            <i class="iconfont">&#xe65d;</i>
          </div>
        </div>
      </div>
      <div class="play-left">
        <i class="iconfont Last" @click="preMusic">&#xe602;</i>
        <i class="iconfont Last" @click="musicPlay()" v-if="musicStatus"
          >&#xe609;</i
        >
        <i class="iconfont control" @click="musicStop()" v-else>&#xe662;</i>
        <i class="iconfont Next" @click="nextMusic">&#xe603;</i>
      </div>
      <div class="play-right">
        <div class="startTime">{{ currentTime }}</div>
        <el-slider
          :disabled="sliderDisabled"
          @change="sliderChange"
          v-model="sliderTimer"
          :show-tooltip="false"
        >
        </el-slider>
        <div class="startTime">{{ duration }}</div>
        <div class="sound">
          <i class="iconfont" @click="soundHandle" v-if="!musicSound"
            >&#xe63b;</i
          >
          <i class="iconfont" @click="soundHandle" v-else>&#xe666;</i>
          <div class="audio">
            <audio
              autoplay
              @play="playLoad()"
              ref="audio"
              @error="audioError()"
              @timeupdate="timeupdate"
              @playing="musicPlaying()"
              @pause="musicPause()"
              @ended="musicEnded()"
              :src="playList.src"
            ></audio>
          </div>
          <el-slider :show-tooltip="false" v-model="volume"> </el-slider>
        </div>
        <div class="order">
          <el-tooltip
            class="item"
            effect="dark"
            content="顺序播放"
            placement="top"
          >
            <i class="iconfont">&#xe802;</i>
          </el-tooltip>
          <!-- <el-tooltip class="item" effect="dark" content="随即播放" placement="top">
          <i class="iconfont">&#xe71f;</i>
        </el-tooltip>
        <el-tooltip class="item" effect="dark" content="单曲循环" placement="top">
          <i class="iconfont">&#xe66d;</i>
        </el-tooltip> -->
        </div>
        <div class="musicList">
          <el-tooltip
            class="item"
            effect="dark"
            content="播放列表"
            placement="top"
          >
            <i
              class="iconfont"
              @click="isShowMusicplaylist = !isShowMusicplaylist"
              >&#xe83e;</i
            >
          </el-tooltip>
        </div>
      </div>
    </div>
    <transition name="Musicplaylist">
      <music-play-list
        v-show="isShowMusicplaylist"
        :currentIndex="currentIndex"
        :musicList="musicList"
      ></music-play-list>
    </transition>

    <music-lyric ref="musicLyric"></music-lyric>
  </div>
</template>

<script>
import { _getSongUrl, _getLyric } from "@/network/song";
import { playlistTool } from "./playlist";
import { formatDate, deepClone } from "@/common/js/tool";
import MusicPlayList from "./MusicPlayList";
// 歌词
import MusicLyric from "./MusicLyric";
export default {
  data() {
    return {
      playList: {
        src: ""
      },
      musicList: [],
      currentIndex: 0,
      // 当前播放位置
      // slider的位置
      sliderTimer: 0,
      // slider是否禁用
      sliderDisabled: true,
      currentTime: "00:00",
      // 音频的长度
      duration: "00:00",
      showTop: false,
      userinfo: {
        name: "",
        image_url: "",
        id: "",
        song: ""
      },
      // 当前是播放还是暂停
      musicStatus: true,
      // 是否禁音
      musicSound: false,
      // 音量
      volume: 50,
      // 播放列表
      isShowMusicplaylist: false,
      // 播放类型
      playbackType: 1
    };
  },

  methods: {
    audioError() {
      // this.$message.error('没有音频')
    },
    // 开始播放
    musicPlaying() {
      this.$store.commit("editMusicPlay", true);
    },
    // 暂停了
    musicPause() {
      this.$store.commit("editMusicPlay", false);
    },
    // 当歌曲加载完成播放时
    playLoad() {
      this.$refs.audio.volume = this.volume / 100;
      this.duration = this.musicList[this.currentIndex].time;
    },
    // 播放
    musicPlay() {
      // 我是把要播放的放进playlist中  跟最初想的不一样了
      // 这里不知道为什么有时候获取不到this.$refs.audio  返回Undefined
      this.$nextTick(() => {
        if (this.playList.src !== "") {
          this.$refs.audio.volume = 0;
          this.$refs.audio.play();
          this.musicStatus = false;
        }
      });
    },
    // 重新播放
    musicLoad() {
      this.$refs.audio.load();
    },
    // 渐变声音
    musicGradients(type) {
      return new Promise((resolve, reject) => {
        let num = 4; // 在对半除的时候 3 4次就没声音了
        for (let i = 1; i <= num; i++) {
          setTimeout(() => {
            if (i == num) resolve();
            if (type === "down") {
              this.$refs.audio.volume = this.$refs.audio.volume / 2;
            } else {
              // 这里不知道怎么做到播放也是渐变的。。。
            }
          }, i * 100);
        }
      });
    },
    // 暂停
    async musicStop() {
      this.musicStatus = true;
      await this.musicGradients("down");
      this.$refs.audio.pause();
    },
    // 根据获取到的数据 发起获取URL请求 然后对返回值过滤 添加到播放中
    getSongUrl(id) {
      // _getSongUrl(id).then(result => {
      //   const res = new playlistTool(result.data[0]);
      //   if (!res.src) {
      //     this.nextMusic();
      //     return this.$message.info("需要付费或其他原因");
      //   }
      //   if (res.src == this.playList.src) {
      // 再判断如果相同 那是不是结束了
      //     if (this.$refs.audio && this.$refs.audio.ended) {
      // 重新播放
      //       return this.musicLoad();
      //     }
      //     return this.$message.info("正在播放中");
      //   }
      //   this.playList = res;
      // console.log(res);
      // 放在这里只有歌真正变化了再去修改
      //   this.$store.commit("editSongDetai", this.musicList[this.currentIndex]);
      //   this.musicStatus = false;
      // });
    },
    getMusicUrl(music) {
      const res = new playlistTool(music);
      if (res.src == this.playList.src) {
        // 再判断如果相同 那是不是结束了
        if (this.$refs.audio && this.$refs.audio.ended) {
          // 重新播放
          return this.musicLoad();
        }
        return this.$message.info("正在播放中");
      }
      this.playList = res;
      this.$store.commit("editSongDetai", this.musicList[this.currentIndex]);
      this.musicStatus = false;
    },
    // 歌曲时间信息
    timeupdate() {
      // 在歌曲停止的一瞬间 this.$refs.audio 是undefined 然后就会报错。。
      if (!this.$refs.audio) return;
      // 当前时间
      // console.log(this.$refs.audio);
      // console.log(this.$refs.audio.duration);
      let duration = new Date(this.$refs.audio.duration * 1000);
      let currentTime = new Date(this.$refs.audio.currentTime * 1000);
      this.currentTime = formatDate(currentTime, "mm:ss");
      // 滑块的位置
      this.sliderTimer = (currentTime / duration) * 100;

      this.$store.commit("editCurrentTime", currentTime);
      // console.log(this.sliderTimer);
    },
    // 更新位置
    sliderChange(val) {
      // 更新的时间位置等于  新位置的百分比 * 总的时间长度
      this.$refs.audio.currentTime = (val / 100) * this.$refs.audio.duration;
    },
    // 下一首
    nextMusic() {
      // console.log('next');
      // 是不是最后一首
      if (this.musicList.length - 1 == this.currentIndex) {
        // 判断数组中是否还有歌 如果有 就跳到第一首 这个就是列表循环播放   我实现不了。。
        // 是否播放完成
        if (this.$refs.audio.ended) {
          // 播放完成又点击 重新播放
          return this.musicLoad();
        } else {
          return;
        }
        this.musicStatus = true;
        return this.$message.info("正在播放中");
      }
      this.currentIndex += 1;
      this.init();
    },
    // 上一首
    preMusic() {
      if (this.currentIndex === 0) return this.$message.info("这是第一首了");
      this.currentIndex -= 1;
      this.init();
    },
    // 播放结束的时候
    musicEnded() {
      // 最后一首就不自动播放了
      if (this.musicList.length - 1 == this.currentIndex) {
        // 那就吧左下角隐藏
        return (this.showTop = false);
      }
      this.nextMusic();
    },
    // 显示左下角歌曲信息
    showTopHandle(info) {
      this.userinfo = info;
      this.showTop = true;
    },
    // 控制禁音
    soundHandle() {
      this.musicSound = !this.musicSound;
      if (this.musicSound) {
        this.volume = 0;
      } else {
        this.volume = 50;
      }
    },
    // 重置一下播放器信息
    resetInfo() {
      this.sliderTimer = 0;
      this.playList = {
        src: ""
      };
      this.musicList = [];
      this.currentTime = "00:00";
      this.duration = "00:00";
      this.showTop = false;
    },
    // 播放控制
    init() {
      this.getMusicUrl(this.musicList[this.currentIndex]);
      this.showTopHandle(this.musicList[this.currentIndex]);
    },
    // 显示歌词面板
    showLyricBox() {
      this.$store.commit("editshowLyric", true);
    }
  },
  mounted() {
    this.$bus.$on("playMusic", (index, list) => {
      // 对传递过来的list 克隆一下
      let cloneList = [];
      deepClone(cloneList, list);
      this.currentIndex = index;
      this.musicList = cloneList;
      this.init();
    });
    // todos: 搜索一首歌的时候 播放是把这歌插入到当前播放列表的中 位置是当前播放的位置
    this.$bus.$on("pushPlayMusic", item => {
      if (this.musicList.length === 0) {
        this.musicList.push(item);
        this.init();
      } else {
        // 需要判断歌曲是否在当前数组中
        let index = this.musicList.findIndex(value => {
          return item.id == value.id;
        });
        if (index !== -1) {
          this.currentIndex = index;
          this.init();
          return;
        }
        this.currentIndex++;
        this.musicList.splice(this.currentIndex, 0, item);
        this.init();
      }
    });
  },
  beforeDestroy() {
    // 貌似这个组件会一直存在。。。
    // this.$bus.$off('playMusic')
    // this.$bus.$off('pushPlayMusic')
  },
  watch: {
    // 监听音量
    volume: function(val) {
      if (val == 0) {
        this.musicSound = true;
      } else {
        this.musicSound = false;
      }
      this.$refs.audio.volume = val / 100;
    },
    // 监听是否有数据
    sliderTimer() {
      if (this.musicList.length === 0) {
        this.sliderDisabled = true;
      } else {
        this.sliderDisabled = false;
      }
    },
    // 当state中的数据发生变化 就代表歌词面板变化 左下角反着来
    showLyric: function(val) {
      this.showTop = !val;
    }
  },
  components: {
    MusicPlayList,
    MusicLyric
  },
  computed: {
    // 变成计算属性
    showLyric() {
      return this.$store.state.showLyric;
    }
  }
};
</script>

<style lang="less" scoped>
i {
  cursor: pointer;
}
.Musicplaylist-enter-active,
.Musicplaylist-leave-active {
  transition: opacity 0.3s;
}
.Musicplaylist-enter, .Musicplaylist-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
.musicplay {
  display: flex;
  position: absolute;
  z-index: 150;
  bottom: 0;
  width: 100%;
  border-top: 1px solid #e1e1e2;
  background-color: #f6f6f8;
  height: 50px;
  .play-top {
    position: absolute;
    top: -60px;
    width: 200px;
    height: 60px;
    box-sizing: border-box;
    padding: 5px;
    background-color: #fff;
    display: flex;
    align-items: center;
    .image {
      position: relative;
      width: 45px;
      height: 45px;
      img {
        width: 100%;
        height: 100%;
      }
      .imageIconfont {
        display: none;
        position: absolute;
        left: 0;
        width: 45px;
        color: #fff;
        text-align: center;
        height: 45px;
        font-size: 30px;
        background-color: rgba(0, 0, 0, 0.4);
      }
    }
    .image:hover {
      .imageIconfont {
        display: inline-block;
      }
    }
    .userinfo {
      flex: 1;
      font-size: 12px;
      display: flex;
      justify-content: space-between;
      margin-left: 10px;
      .info {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        .title {
          width: 110px;
          white-space: nowrap;
          text-overflow: ellipsis;
          overflow: hidden;
        }

        .nickname {
          width: 76px;
          white-space: nowrap;
          text-overflow: ellipsis;
          overflow: hidden;
        }
      }

      .warp {
        width: 10px;
        display: flex;
        flex-wrap: wrap;
        margin-right: 5px;
        align-items: center;
        justify-content: center;
        i.iconfont {
          font-size: 10px;
          cursor: pointer;
        }
      }
    }
  }
  .play-left {
    width: 200px;
    align-items: center;
    display: flex;
    justify-content: center;
    i {
      cursor: pointer;
      font-size: 30px;
      color: #c62f2f;
      border-radius: 50%;
      border: 1px solid #fff;
      margin: 10px;
      vertical-align: top;
    }
  }
  .play-right {
    flex: 1;
    display: flex;
    align-items: center;
    .startTime {
      // width: 20px;
      margin-right: 20px;
    }
    .el-slider {
      flex: 1;
      margin-right: 10px;
    }
    .sound {
      display: flex;
      align-items: center;
      width: 130px;
      .el-slider {
        flex: 1;
      }
      i {
        width: 25px;
        margin-right: 5px;
      }
    }
    .startTime {
      margin-left: 20px;
    }
    .order {
      margin: 0 10px;
    }
    .musicList {
      margin: 0 10px;
      margin-right: 30px;
    }
  }
}
</style>

