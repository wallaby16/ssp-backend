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
      this.$store.commit('setNotification', {
        notification: {
          type: res.status === 200 ? 'success' : 'danger',
          message: res.body.message
        }
      });
    }

    if (res.url !== '/login' && res.status === 401) {
      this.$store.commit('setUser', {
        user: null
      });
      this.$store.commit('setNotification', {
        notification: {
          type: 'danger',
          message: 'Dein Token ist abgelaufen. Bitte logge dich neu ein'
        }
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
