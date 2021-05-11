<template>
  <div class="header">
    <div class="logo cursorPointer" @click="$router.push('/home/discover')">
      <img src="@/common/img/logo.jpg" alt="" />
      <p>网易云音乐</p>
      <i class="el-icon-arrow-left cursorPointer" @click="backRouter"></i>
    </div>
    <div class="search">
      <el-input
        ref="searchInput"
        @focus="focusHandle"
        @blur="blurHandle"
        size="small"
        suffix-icon="cursorPointer el-icon-search"
        placeholder="请输入内容"
        v-model="searchValue"
      >
      </el-input>
      <!-- 弹出搜索框 -->

      <search-suggest
        :content="searchValue"
        v-show="searchSuggest"
        :suggestList="suggestList"
      ></search-suggest>
    </div>
    <div class="content">
      <div @click="showUserFormDialog" class="login cursorPointer">
        <el-avatar size="small" :src="loginImgSrc"></el-avatar>
        <span>{{ loginState }}</span>
      </div>
      <div class="map">
        <i class="el-icon-minus cursorPointer"></i>
        <i class="el-icon-document-copy cursorPointer"></i>
        <i class="el-icon-close cursorPointer"></i>
      </div>
    </div>

    <!-- 用户对话框 -->
    <el-dialog
      :modal="false"
      class="loginDialog"
      title=""
      :visible.sync="userFormDialogVisible"
      width="25%"
      @close="userFormClose"
    >
      <img src="@/common/img/logo.jpg" alt="" />
      <br />
      <div style="width: 100%">
        <el-tabs v-model="activeName">
          <el-tab-pane label="登录" name="first">
            <el-form
              :model="userForm"
              :rules="userFormRules"
              ref="userFormRef"
              class="demo-ruleForm"
            >
              <el-form-item prop="username">
                <el-input
                  placeholder="请输入您的邮箱"
                  size="small"
                  v-model="userForm.username"
                ></el-input>
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  placeholder="请输入您的密码"
                  type="password"
                  size="small"
                  v-model="userForm.password"
                ></el-input>
              </el-form-item>
            </el-form>
            <button class="loginBtn" @click="userLogin">登录</button>
          </el-tab-pane>
          <el-tab-pane label="注册" name="second">
            <el-form
              :model="userRegisterForm"
              :rules="userRegisterFormRules"
              ref="userRegisterFormRef"
              class="demo-ruleForm"
            >
              <el-form-item prop="email">
                邮箱：<el-input
                  placeholder="请输入您的邮箱"
                  size="small"
                  v-model="userRegisterForm.email"
                ></el-input>
              </el-form-item>
              <el-form-item prop="username">
                用户名：<el-input
                  placeholder="请输入您的邮箱"
                  size="small"
                  v-model="userRegisterForm.username"
                ></el-input>
              </el-form-item>
              <el-form-item prop="password">
                密码：<el-input
                  placeholder="请输入您的密码"
                  type="password"
                  size="small"
                  v-model="userRegisterForm.password"
                ></el-input>
              </el-form-item>
              <el-form-item prop="password">
                验证码：<el-input
                  placeholder="请输入您的密码"
                  type="password"
                  size="small"
                  v-model="userRegisterForm.verification"
                ></el-input>
                <el-button
                  type="primary"
                  v-show="showIt"
                  @click="sendVerification"
                  style="width: 150px;"
                  >获取验证码</el-button
                >
                <el-button
                  type="primary"
                  v-show="!showIt"
                  @click="sendVerification"
                  style="width: 150px;"
                  disabled
                  >{{ btntxt }}</el-button
                >
              </el-form-item>
            </el-form>
            <button class="registerBtn" @click="userRegister">注册</button>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import searchBox from "./searchBox";
import SearchSuggest from "./SearchSuggest";
import { _getSearchHot, _getSearchSuggest } from "@/network/search";
import { _userLogin, _userRegister, _userVerification } from "@/network/user";
export default {
  data() {
    // 手机号码验证规则
    var checkMobile = (rule, value, cb) => {
      const regMobile = /^\w{3,15}\@\w+\.[a-z]{2,3}$/;
      if (regMobile.test(value)) {
        return cb();
      }
      cb(new Error("请输入合法的邮箱"));
    };
    return {
      showIt: true,
      btntxt: "获取验证码",
      time: 0,
      activeName: "first",
      readme: false,
      // 搜索内容
      searchValue: "",
      // 登录用户名文本
      loginState: "未登录",
      // 图片地址
      loginImgSrc:
        "https://cube.elemecdn.com/9/c2/f0ee8a3c7c9638a54940382568c9dpng.png",
      userFormDialogVisible: false,
      // 用户信息
      userForm: {
        username: "",
        password: ""
      },
      userRegisterForm: {
        email: "",
        username: "",
        verification: "",
        password: ""
      },
      // 用户验证规则
      userFormRules: {
        username: [
          { required: true, message: "请输入您的邮箱", trigger: "blur" },
          { validator: checkMobile, trigger: "blur" }
        ],
        password: [{ required: true, message: "请输入密码", trigger: "blur" }]
      },
      userRegisterFormRules: {
        email: [
          { required: true, message: "请输入您的邮箱", trigger: "blur" },
          { validator: checkMobile, trigger: "blur" }
        ],
        username: [{ required: true, message: "请输入密码", trigger: "blur" }],
        verification: [
          { required: true, message: "请输入密码", trigger: "blur" }
        ],
        password: [{ required: true, message: "请输入密码", trigger: "blur" }]
      },
      // 输入框的显示隐藏
      focusFlag: false,
      // 热搜数据
      searchHotList: [],
      // 搜索建议框
      searchSuggest: false,
      // 定时器需要放置在这里
      timer: null,
      // 搜索建议对象
      suggestList: {}
    };
  },
  methods: {
    // 关闭对话框的时候
    userFormClose() {
      this.$refs.userFormRef.resetFields();
    },

    getUserInfo(token) {
      this.$http
        .get("/user/me", {
          headers: {
            "sta-token": token
          }
        })
        .then(res => {
          this.loginImgSrc = res.data["image_url"];
          this.loginState = res.data["name"];
          this.userFormDialogVisible = false;
          this.$store.commit("addUser", res.data);
        });
    },
    getUserPlaylist() {},
    // 多少秒后重新发送
    makeTimer() {
      if (this.time > 0) {
        this.showIt = false;
        this.time--;
        this.btntxt = this.time + "s重新获取";
        setTimeout(this.makeTimer, 1000);
      } else {
        this.showIt = true;
        this.time = 0;
        this.btntxt = "获取验证码";
        this.disabled = false;
      }
    },
    sendVerification() {
      // _userVerificatio
      //60s之后才能点击
      this.time = 60;
      this.makeTimer();
      _userVerification(this.userRegisterForm.email).then(result => {
        if (result.code !== 0) {
          return this.$message.error(result.message);
        }
        return this.$message.sucess(
          "发送成功，请留意邮箱，若未收到邮件请查看垃圾邮件"
        );
      });
    },

    // 登录功能
    userLogin() {
      let tokenTmp = "";
      this.$refs.userFormRef.validate(async item => {
        if (!item) return;
        _userLogin({
          username: this.userForm.username,
          password: this.userForm.password,
          readme: this.readme
        }).then(result => {
          if (result.code !== 0) {
            return this.$message.error(result.message);
          }
          this.$message.success("登录成功");
          tokenTmp = result.data[this.$tokenStr];
          localStorage.setItem(this.$tokenStr, result.data[this.$tokenStr]);
          console.log(localStorage.getItem(this.$tokenStr));
          this.$store.commit("storeToken", tokenTmp);
          this.$http.defaults.headers.common[this.$tokenStr] = tokenTmp;
          this.getUserInfo(tokenTmp);
        });
      });
    },
    // 用户注册
    userRegister() {
      let tokenTmp = "";
      this.$refs.userRegisterFormRef.validate(async item => {
        if (!item) return;
        _userRegister({
          user: {
            username: this.userRegisterForm.email,
            password: this.userRegisterForm.password,
            name: this.userRegisterForm.username
          },
          code: this.userRegisterForm.verification
        }).then(result => {
          if (result.code !== 0) {
            return this.$message.error(result.message);
          }
          this.$message.success("注册成功");
          tokenTmp = result.data[this.$tokenStr];
          localStorage.setItem(this.$tokenStr, result.data[this.$tokenStr]);
          console.log(localStorage.getItem(this.$tokenStr));
          this.$store.commit("storeToken", tokenTmp);
          this.$http.defaults.headers.common[this.$tokenStr] = tokenTmp;
          this.getUserInfo(tokenTmp);
        });
      });
    },
    // 后退
    backRouter() {
      this.$router.go(0);
      // console.log(this.$route.path);
    },
    // 点击搜索框
    focusHandle() {
      // 判断内容
      // console.log('dian');
      if (this.searchValue.trim().length === 0) {
        this.focusFlag = true;
        this.searchSuggest = false;
      } else {
        this.searchSuggest = true;
        this.focusFlag = false;
      }
    },
    // 失去焦点
    blurHandle() {
      // 这里是因为当你要点击里面的内容的时候
      //  结果失去焦点后如果直接隐藏 那你就点不到了
      setTimeout(() => {
        this.focusFlag = false;
        this.searchSuggest = false;
      }, 100);
    },
    // 显示登录用户对话框
    showUserFormDialog() {
      if (this.loginState === "未登录") {
        this.userFormDialogVisible = true;
      } else {
        // todos:实现退出登录功能
        // 清除图片地址 和文本
      }
    },
    // 获取热搜数据
    getSearchHot() {
      _getSearchHot().then(result => {
        console.log(result);
        this.searchHotList = result.data;
      });
    },
    // 获取搜索建议
    getSearchSuggest(keywords) {
      return new Promise((resolve, reject) => {
        _getSearchSuggest(keywords).then(res => {
          this.suggestList = res.data;
          console.log(this.suggestList);
          resolve();
        });
      });
    }
  },
  components: {
    searchBox,
    SearchSuggest
  },
  created() {
    this.getSearchHot();
  },
  watch: {
    searchValue: function(val) {
      if (val.trim() !== "") {
        // 不为空
        clearTimeout(this.timer);
        this.timer = setTimeout(async () => {
          await this.getSearchSuggest(val);
          if (!this.suggestList) return;
          this.focusFlag = false;
          this.searchSuggest = true;
        }, 250);
      } else {
        // 内容为空
        clearTimeout(this.timer);
        this.searchSuggest = false;
        this.focusFlag = true;
      }
    }
  }
};
</script>

<style lang="less">
.header {
  // z-index: 150;
  box-sizing: border-box;
  padding: 0 10px;
  width: 100%;
  display: flex;
  height: 50px;
  background-color: #c62f2f;
  align-items: center;
  .logo {
    width: 190px;
    display: flex;
    align-items: center;
    font-size: 14px;
    font-weight: 700;
    color: #fff;
    img {
      width: 30px;
      height: 30px;
      border-radius: 10px;
      margin-right: 10px;
    }
    i {
      flex: 1;
      text-align: right;
      margin-right: 10px;
      font-size: 20px;
    }
  }
  .search {
    display: flex;
    align-items: center;
    margin-left: 50px;
    position: relative;
    .el-icon-search:before {
      font-size: 20px;
    }
  }
  .content {
    flex: 1;
    display: flex;
    justify-content: flex-end;
    .login {
      display: flex;
      align-items: center;
      span {
        margin-left: 10px;
        color: #fff;
      }
    }
    .map {
      display: flex;
      align-items: center;
      margin-left: 50px;
      i {
        margin: 0 10px;
        color: #fff;
      }
    }
  }
}
.loginDialog {
  .el-tabs__content {
    width: 100%;
  }
  img {
    margin: 0 auto;
    display: block;
    width: 80px;
    border-radius: 10px;
  }
  .demo-ruleForm {
    margin: 50px 0 100px;
    padding: 0 50px;
  }

  .registerBtn {
    margin: 0 auto;
    display: block;
    width: 220px;
    height: 45px;
    background-color: #00d0ff;
    border: 0;
    outline: none;
    cursor: pointer;
    font-size: 18px;
    color: #fff;
    font-weight: 600;
    border-radius: 10px;
  }
  .loginBtn {
    margin: 0 auto;
    display: block;
    width: 220px;
    height: 45px;
    background-color: #ff0000;
    border: 0;
    outline: none;
    cursor: pointer;
    font-size: 18px;
    color: #fff;
    font-weight: 600;
    border-radius: 10px;
  }
}
</style>
