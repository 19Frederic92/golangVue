// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import Users from '../components/Users.vue'
import Entite from '../components/Entite.vue'
import tokenDetail from '../views/tokenDetail.vue'
import EntiteDetail from '../views/EntiteDetail.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
  /*  {
      path: '/',
      name: 'home',
      component: Users
    }, */

    {
      path: '/',
      name: 'home',
      component: Entite
    },
    {
      path: '/token/:id',
      name: 'tokenDetail',
      component: tokenDetail,
      props: true // Active le passage des paramètres comme props
    },
    {
      path: '/entites/:id',
      name: 'entiteDetail',
      component: EntiteDetail,
      props: true // Active le passage des paramètres comme props
    }

  ]
})

export default router
