import { createRouter, createWebHistory } from 'vue-router'
import HomePage from './views/HomePage.vue'
import TokenDetailPage from './views/TokenDetailPage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage
  },
  {
    path: '/token/:id',
    name: 'TokenDetail',
    component: TokenDetailPage,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

export default router
