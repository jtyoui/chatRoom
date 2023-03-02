<template>
  <div class="box">
    <el-container style="height: 800px; border: 1px solid #eee">
      <el-aside style="background-color: rgb(238, 241, 246)" width="200px">
        <h1>在线人数：{{ number }}</h1>

        <div class="infinite-list-wrapper" style="overflow:auto;border-top: 1px solid #0072C1">
          <div v-for="user of userList" style="list-style-type: none">
            <span style="cursor:pointer;font-size: 16px;margin: auto" @click="userCopyChat(user)">{{ user }}</span>
          </div>
        </div>
      </el-aside>

      <el-container>
        <el-header>
          <span style="margin-left: 400px; font-size: 16px;color: rgba(51,33,210,0.42)">欢迎进入聊天室</span>
          <span style="float: right; font-size: 14px">当前用户：{{ username }}</span>
        </el-header>

        <el-main id="scrollText">
          <div v-for="msg of data" class="chat">

            <div v-if="msg.right" class="right cams">
              <img alt="" class="headIcon radius" src="../static/img/B.jpg">
              <span class="name"><p class="title">{{ msg.time }}</p></span>
              <span class="content"> {{ msg.right }} </span>
            </div>

            <div v-for="message of msg.left">
              <div v-if="message.content" class="left cams">
                <img alt="" class="headIcon radius" src="../static/img/A.jpg"/>
                <span class="name"><p class="title">{{ message.user }}-{{ message.time }}</p></span>
                <span class="content">{{ message.content }}</span>
              </div>
            </div>
          </div>
        </el-main>

        <el-input v-model="send" :rows="2" placeholder="请输入内容" type="textarea"
                  @keydown.enter="handleSendClick"></el-input>
      </el-container>
    </el-container>
  </div>
</template>

<script lang="ts">
import {ref, reactive, onMounted} from "vue";
import {useRouter} from "vue-router";

export default {
  name: 'Home',
  setup() {
    const ip = import.meta.env["VITE_WS_IP"]
    const data = reactive([{
      left: [{user: '', content: '', time: ''}],
      right: '',
      time: '',
    }])
    let username = ref('')
    const route = useRouter()
    const send = ref('')
    const ws = new WebSocket(ip)
    let request = {type: '', content: ''}
    const number = ref(0)
    const userList = ref([])

    onMounted(() => {
      const _username = localStorage.getItem("username")
      const _data = localStorage.getItem("chat")
      if (_username && _username.trim().length) {
        username.value = _username
      } else {
        route.push('/login')
      }
      if (_data) {
        JSON.parse(_data).forEach((e: { left: { user: string; content: string; time: string; }[]; right: string; time: string; logout: string[]; login: string[]; }) => {
          data.push(e)
        })
      }
    })

    setInterval(() => {
      const textarea = document.getElementById('scrollText');
      if (textarea) {
        textarea.scrollTop = textarea.scrollHeight;
      }
    }, 100)

    ws.onmessage = function (e) {
      const msg = JSON.parse(e.data)
      switch (msg.type) {
        case 'handshake':
          request['type'] = 'login'
          request['content'] = username.value
          ws.send(JSON.stringify(request))
          break
        case 'user':
          if (msg.user === username.value) {
            data.push({
              left: [],
              right: msg.content,
              time: msg.time,
            })
          } else {
            data[data.length - 1].left.push({user: msg.user, content: msg.content, time: msg.time})
          }
          localStorage.setItem("chat", JSON.stringify(data))
          break
        case 'login':
          break
        case 'logout':
          break
      }
      if (msg.user_list) {
        userList.value = msg.user_list
        number.value = userList.value.length
      }
    }

    const handleSendClick = () => {
      if (send.value.trim().length) {
        request['type'] = 'user'
        request['content'] = send.value
        ws.send(JSON.stringify(request))
      }
      send.value = ''
    }

    const userCopyChat = (user: string) => {
      send.value = ''
      send.value = user + "@"
    }

    ws.onerror = function (e) {
      alert("聊天异常，请刷新浏览器！")
    }

    return {handleSendClick, send, data, username, number, userList, userCopyChat}
  }
}
</script>
