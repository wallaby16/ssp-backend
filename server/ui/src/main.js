import Vue from 'vue';
import VueRouter from 'vue-router';
import Buefy from 'buefy';
import VueResource from 'vue-resource';
// Styles
import 'buefy/lib/buefy.css';
import './theme.css';
// Store
import store from './store';
// Components
import {GlobalComponents, LocalComponents} from './components';
// Router
import router from './router';

// Mixins
Vue.use(VueRouter);
Vue.use(Buefy);
Vue.use(VueResource);

new Vue({
  router,
  store,
  components: LocalComponents,
  el: '#app',
  render: h => h(GlobalComponents.App)
});
