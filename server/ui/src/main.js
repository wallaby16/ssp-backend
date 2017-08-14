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

// Http interceptors: Global response handler
Vue.http.interceptors.push(function(request, next) {
  next(function(res) {
    if (res.body.message) {
      this.$toast.open({
        type: res.status === 200 ? 'is-success' : 'is-danger',
        message: res.body.message,
        duration: 5000
      });

    }
  });
});

// Http interceptors: Add Auth-Header if token present
Vue.http.interceptors.push(function(request, next) {
  if (store.state.user) {
    request.headers.set('Authorization', `Bearer ${store.state.user.token}`);
  }
  next();
});

new Vue({
  router,
  store,
  components: LocalComponents,
  el: '#app',
  render: h => h(GlobalComponents.App)
});
