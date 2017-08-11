import Vue from 'vue'
import VueRouter from 'vue-router'
import Buefy from 'buefy'
import VueResource from 'vue-resource'

import App from './App.vue'
import Home from './components/Home.vue'
import Nav from './components/Nav.vue'

import EditQuota from './components/ose/EditQuota.vue'

Vue.use(VueRouter)
Vue.use(Buefy)
Vue.use(VueResource)

// Components
Vue.component('navbar', Nav)
Vue.component('editquota', EditQuota)

// Routing
const routes = [
  {
    path: '/', component: Home },
  {
    path: '/ose/editquotas', component: EditQuota
  }
]

const router = new VueRouter({routes})

new Vue({
  router,
  el: '#app',
  render: h => h(App)
})
