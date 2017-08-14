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

// Handle http errors globally
Vue.http.interceptors.push(function(request, next) {
  next(function(res) {
    if (res.status !== 200 && res.body.message) {
      this.$toast.open({
        type: 'is-danger',
        message: res.body.message,
        duration: 5000
      });
    }
  });
});

new Vue({
  router,
  store,
  components: LocalComponents,
  el: '#app',
  render: h => h(GlobalComponents.App)
});
