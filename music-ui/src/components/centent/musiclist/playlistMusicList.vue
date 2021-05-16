<template>
  <div
    class="TableList"
    :infinite-scroll-immediate="false"
    v-infinite-scroll="tableScroll"
  >
    <!-- 这里不知道为什么加类名没有用  只能一条一条加了。 -->
    <el-table
      :show-header="hiddenTableHeader"
      stripe
      :header-cell-style="{ padding: '2px 0', fontSize: '12px' }"
      :row-style="{
        padding: '2px 0',
        fontsize: '12px',
        backgroundColor: '#f5f5f7'
      }"
      :cell-style="{
        padding: '2px 0',
        fontSize: '12px',
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis'
      }"
      :row-class-name="tableRowClassName"
      :data="tracklist"
      style="width: 100%"
    >
      <el-table-column type="index" v-if="hiddenIndex"> </el-table-column>
      <el-table-column v-if="hiddenOperation" label="操作" width="70px">
        <template slot-scope="scope">
          <i class="iconfont" @click="deleteMusicForPlaylist(scope.$index)"
            >&#xe606;</i
          >
          <i class="iconfont" @click="rowClick(scope.$index)">&#xe6cf;</i>
        </template>
      </el-table-column>
      <el-table-column label="音乐名" prop="name" :show-overflow-tooltip="true">
      </el-table-column>
      <el-table-column
        v-if="hiddenSinger"
        label="歌手"
        prop="creator_names"
        :show-overflow-tooltip="true"
      >
      </el-table-column>
      <el-table-column
        v-if="hiddenAlbum"
        label="标题"
        :show-overflow-tooltip="true"
        width="380px"
        prop="title"
      >
      </el-table-column>
      <el-table-column
        label="时长"
        prop="play_time"
        :formatter="parsePlayTime"
        width="70px"
      >
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { formatDate } from "../../../common/js/tool.js";
import { MessageBox } from "element-ui";
import { _deleteMusicToPlaylist } from "@/network/discover/discover";
export default {
  data() {
    return {};
  },
  props: {
    tracklist: {
      type: Array,
      default() {
        return [];
      }
    },
    hiddenTableHeader: {
      type: Boolean,
      default() {
        return true;
      }
    },
    hiddenSinger: {
      type: Boolean,
      default() {
        return true;
      }
    },
    hiddenAlbum: {
      type: Boolean,
      default() {
        return true;
      }
    },
    hiddenIndex: {
      type: Boolean,
      default() {
        return true;
      }
    },
    hiddenOperation: {
      type: Boolean,
      default() {
        return true;
      }
    },
    currentIndex: {
      type: Number
    }
  },
  filters: {
    formatDate(time) {
      var data = new Date(time);
      return formatDate(data, "mm:ss");
    }
  },
  methods: {
    tableRowClassName({ row, rowIndex }) {
      //把每一行的索引放进row
      row.index = rowIndex;
      // 这里设置的样式不知道为什么不生效，。
      // if (rowIndex % 2) {
      //   return 'el-row-gehangbianse'
      // }
      // 点击的那一个变色pink
      if (this.currentIndex == rowIndex) {
        return "pink";
      } else {
        return "";
      }
    },
    deleteMusicForPlaylist(mindex) {
      console.log("key = ", this.$route.query.id);
      var pidStr = this.$route.query.id;
      var pid = parseInt(pidStr);
      var msg = "是否将" + this.tracklist[mindex].name + "移除此歌单?";
      var mid = this.tracklist[mindex].id;
      MessageBox.confirm(msg, "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          _deleteMusicToPlaylist({
            pid: pid,
            mid: mid
          }).then(res => {
            if (res.code !== 0) {
              this.$message.error(res.message);
              return;
            }
            this.tracklist.splice(mindex, 1);
            this.$message({
              type: "success",
              message: "删除成功!"
            });
          });
        })
        .catch(() => {});
    },
    parsePlayTime(tm) {
      var data = new Date(tm.play_time);
      return formatDate(data, "mm:ss");
    },
    // 点击了某一行
    rowClick(row) {
      console.log("row = ", row);
      // console.log(123);
      // 交给父组件处理
      this.$emit("rowClick", row, this.tracklist);
      // this.$bus.$emit('playMusic', row.index, this.tracklist)
    },
    // 到底
    tableScroll() {
      this.$emit("tableScroll");
    }
  }
};
</script>

<style lang="less" scoped>
.TableList {
  tr {
    height: 20px !important;
    i {
      margin: 0 2px;
      font-size: 14px;
    }
    i:hover {
      color: #000;
    }
  }
  tr:hover {
    td {
      div {
        color: #000;
      }
    }
  }
}
.el-table tr:nth-child(2) {
  background-color: red;
}
</style>
