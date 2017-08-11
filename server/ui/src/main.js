import Vue from 'vue';
import Vuex from 'vuex';
import VueRouter from 'vue-router';
import Buefy from 'buefy';
import VueResource from 'vue-resource';

// Styles
import 'buefy/lib/buefy.css';
import './theme/bulmaswatch.min.css';

// Global Components
import App from './components/App.vue';
import Home from './components/shared/Home.vue';
import Nav from './components/shared/Nav.vue';
import Login from './components/shared/Login.vue';

// OSE-Components
import EditQuota from './components/ose/EditQuota.vue';
import NewProject from './components/ose/NewProject.vue';
import NewTestProject from './components/ose/NewTestProject.vue';
import UpdateBilling from './components/ose/UpdateBilling.vue';
import NewServiceAccount from './components/ose/NewServiceAccount.vue';

// Store
import store from './store'

// Mixins
Vue.use(VueRouter);
Vue.use(Buefy);
Vue.use(VueResource);

// Components
Vue.component('login', Login);
Vue.component('navbar', Nav);

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
  store,
  components: {
    EditQuota,
    NewTestProject,
    NewProject,
    NewServiceAccount,
    UpdateBilling
  },
  el: '#app',
  render: h => h(App)
});
