<template>
  <div class="wrap">
    <el-form :inline="true" class="demo-form-inline">
      <el-form-item label="用户名：">
        <el-input v-model="username" placeholder="请输入用户名"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleEnterClick">进入聊天室</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
import {ref, onMounted} from "vue";
import {useRouter} from "vue-router";

export default {
  name: 'Login',
  setup() {
    const username = ref("")
    const route = useRouter()

    onMounted(() => {
      const _username = localStorage.getItem("username")
      if (_username) {
        route.push('/')
        return
      }
    })
    const handleEnterClick = () => {
      const _username = username.value.trim()
      if (_username.length < 2) {
        alert("用户名长度不能小于2位")
      } else {
        localStorage.setItem("username", _username)
        username.value = ''
        route.push('/')
      }
    }
    return {username, handleEnterClick}
  }
}
</script>

<style>
.wrap {
  line-height: 200px; /*垂直居中关键*/
  text-align: center;
  height: 200px;
  font-size: 36px;
  background-color: #5e4242;
}
</style>
