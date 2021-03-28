<template>
  <input
    type="button"
    value="获取我信息"
    name="info"
    id="info"
    @click="getMeInfo"
  />
</template>

<script>
export default {
  methods: {
    getMeInfo () {
      const token = localStorage.getItem(this.$tokenStr)
      console.log(token)
      this.$http.get('/user/me').then(resp => {
        console.log(resp)
      })
    }
  }
}
</script>

<style lang="less" scoped></style>
