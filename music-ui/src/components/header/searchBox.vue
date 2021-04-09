<template>
  <div class="searchBox">
    <!-- <transition name="slide-fade"> -->
    <div class="searchlist" style="overflow-y:auto;">
      <div class="title">热搜榜</div>
      <div class="tablelist">
        <table cellpadding="0">
          <tr
            @click="search(item.searchWord)"
            v-for="(item, index) in searchHotList"
            :key="index"
          >
            <td class="index">{{ index + 1 }}</td>
            <td>
              <div>
                <div class="info">
                  <span class="searchWord">{{ item.searchWord }}</span>
                  <span class="score">{{ item.score }}</span>
                </div>
                <div class="contetn">{{ item.content }}</div>
              </div>
            </td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    searchHotList: {
      type: Array
    }
  },
  methods: {
    // 点击搜索
    search(keywords) {
      this.$router.push({
        path: "/home/searchlist",
        query: {
          content: keywords
        }
      });
    }
  }
};
</script>

<style lang="less" scoped>
.searchBox {
  position: relative;
  width: auto;
  height: auto;
  position: absolute;
  left: 0;
  top: 41px;
  z-index: 110;
  .searchlist {
    width: 320px;
    height: 400px;
    background-color: #fff;
    box-shadow: 4px 4px 20px rgba(0, 0, 0, 0.3);
    .title {
      height: 40px;
      line-height: 40px;
      font-size: 15px;
      padding-left: 20px;
    }
    .tablelist {
      table {
        width: 100%;
        border: none;
        border-spacing: 0;
        padding: 0;
        tr {
          height: 55px;
          padding: 0;
          .index {
            box-sizing: border-box;
            width: 50px;
            padding-left: 20px;
            line-height: 50px;
          }
          .info {
            .searchWord {
              color: #333333;
            }
            .score {
              margin-left: 15px;
              color: #cdcdcd;
            }
          }
          .contetn {
            color: #999;
          }
        }
        tr:hover {
          background-color: #ededed;
        }
        tr:nth-child(-n + 3) {
          td:nth-child(1) {
            color: #ff3d3d;
            font-weight: 600;
          }
          td:nth-child(2) {
            .info {
              .searchWord {
                color: #222;
                font-weight: 600;
                font-size: 13px;
              }
            }
          }
        }
      }
    }
  }
}
.searchBox ::after {
  content: "";
  top: -14px;
  left: 14px;
  position: absolute;
  z-index: 100;
  border: 7px solid transparent;
  border-bottom: 7px solid #fff;
}
</style>

