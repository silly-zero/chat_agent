import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// 配置axios
axios.defaults.baseURL = '/api'
axios.defaults.headers.post['Content-Type'] = 'application/json'

const app = createApp(App)

app.use(router)

// 全局挂载axios
app.config.globalProperties.$axios = axios

app.mount('#app')
