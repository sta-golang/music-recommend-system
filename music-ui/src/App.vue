<template>
  <div id="app">
    <the-header />
    <router-view></router-view>
  </div>
</template>

<script>
import header from './components/header.vue'
export default {
  components: { header },
  name: 'app'
}
</script>

<style></style>
