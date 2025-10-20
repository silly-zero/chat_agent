import { createRouter, createWebHistory } from 'vue-router'
import StarSelection from '../views/StarSelection.vue'
import ChatPage from '../views/ChatPage.vue'

const routes = [
  {
    path: '/',
    name: 'StarSelection',
    component: StarSelection
  },
  {
    path: '/chat/:starId',
    name: 'ChatPage',
    component: ChatPage,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router