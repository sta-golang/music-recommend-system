<template>
  <div>
    <h1>扁平响应login登录页面带弹出</h1>

    <div class="w3layoutscontaineragileits">
      <h2>在这里登录</h2>
      <form>
        <input
          id="usernameVal"
          type="email"
          name="Username"
          placeholder="EMAIL"
          required
        />
        <input
          id="passwordVal"
          type="password"
          name="Password"
          placeholder="PASSWORD"
          required
        />
        <ul class="agileinfotickwthree">
          <li>
            <input type="checkbox" id="brand1" value="true" />
            <label for="brand1"><span></span>记住我</label>
            <a href="#">忘记密码?</a>
          </li>
        </ul>
        <div class="aitssendbuttonw3ls">
          <input type="button" value="登录" @click="loginMethod" />
          <p>
            注册一个新账户 <span>→</span>
            <a class="w3_play_icon1" href=""> 点击这里</a>
          </p>
          <div class="clear"></div>
        </div>
      </form>
    </div>

    <!-- for register popup -->

    <div class="w3footeragile">
      <p>
        &copy; This is music-recomment-system project for XUPT-STA
        author(qq:63237777)<br />
        see: <a href="https://github.com/sta-golang">github</a>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    async loginMethod () {
      var formData = {
        username: document.getElementById('usernameVal').value,
        password: document.getElementById('passwordVal').value,
        readme: Boolean(document.getElementById('brand1').checked)
      }
      const res = await this.$http.post('/user/login', formData)
      if (res.data.code !== 0) {
        return this.$message.error(res.data.message)
      }
      localStorage.setItem(this.$tokenStr, res.data.data[this.$tokenStr])
      return this.$message.success('登录成功')
    }
  }
}
</script>

<style lang="less" scoped>
@import '../assets/css/login.css';
@import '../assets/css/popuo-box.css';
.body_back {
  background: url('../assets/images/b2.jpg');
}
</style>
