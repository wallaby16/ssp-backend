import Vue from 'vue';
import VueRouter from 'vue-router';
import Buefy from 'buefy';
import VueResource from 'vue-resource';
import 'buefy/lib/buefy.css';
import './theme/bulmaswatch.min.css';

import App from './App.vue';
import Home from './components/Home.vue';
import Nav from './components/Nav.vue';
import Login from './components/Login.vue';

import EditQuota from './components/ose/EditQuota.vue';
import NewProject from './components/ose/NewProject.vue';
import NewTestProject from './components/ose/NewTestProject.vue';
import UpdateBilling from './components/ose/UpdateBilling.vue';
import NewServiceAccount from './components/ose/NewServiceAccount.vue';

Vue.use(VueRouter);
Vue.use(Buefy);
Vue.use(VueResource);

// Components
Vue.component('login', Login);
Vue.component('navbar', Nav);
Vue.component('editquota', EditQuota);
Vue.component('newproject', NewProject);
Vue.component('newtestproject', NewTestProject);
Vue.component('updatebilling', UpdateBilling);
Vue.component('newserviceaccount', NewServiceAccount);

// Routing
const routes = [
  {
    path: '/', component: Home
  },
  {
    path: '/login', component: Login
  },
  {
    path: '/ose/editquotas', component: EditQuota
  },
  {
    path: '/ose/newtestproject', component: NewTestProject
  },
  {
    path: '/ose/newproject', component: NewProject
  },
  {
    path: '/ose/newserviceaccount', component: NewServiceAccount
  },
  {
    path: '/ose/updatebilling', component: UpdateBilling
  }
];

const router = new VueRouter({routes});

new Vue({
  router,
  el: '#app',
  render: h => h(App)
});
