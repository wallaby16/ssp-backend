import {GlobalComponents, LocalComponents} from './components';
import VueRouter from 'vue-router';
import store from './store/index';

const routes = [
  {
    path: '/', component: LocalComponents.Home
  },
  {
    path: '/login', component: GlobalComponents.Login
  },
  {
    path: '/ose/editquotas', component: LocalComponents.EditQuota
  },
  {
    path: '/ose/newtestproject', component: LocalComponents.NewTestProject
  },
  {
    path: '/ose/newproject', component: LocalComponents.NewProject
  },
  {
    path: '/ose/adminlist', component: LocalComponents.AdminList
  },
  {
    path: '/ose/newserviceaccount', component: LocalComponents.NewServiceAccount
  },
  {
    path: '/ose/updatebilling', component: LocalComponents.UpdateBilling
  },
  {
    path: '/gluster/newvolume', component: LocalComponents.NewVolume
  },
  {
    path: '/gluster/fixvolume', component: LocalComponents.FixVolume
  },
  {
    path: '/gluster/growvolume', component: LocalComponents.GrowVolume
  },
  {
    path: '/ddc/billing', component: LocalComponents.DDCBilling
  },
  {
    path: '/aws/lists3buckets', component: LocalComponents.ListS3Buckets
  },
  {
    path: '/aws/news3bucket', component: LocalComponents.NewS3Bucket
  }
];

const router = new VueRouter({routes});

router.beforeEach((to, from, next) => {
  // Cleanup old notifications
  store.commit('setNotification', {notification: {}});

  // Auth-Protection
  if (to.path == '/login') {
    // Login page is always allowed
    next();
  }
  if (!store.state.user) {
    console.error('Not yet logged in, navigating to login');
    next({path: '/login'});
  } else {
    // Check if token is still valid
    if (store.state.user && store.state.user.exp < Date.now() / 1000) {
      console.error('Token is no longer valid, navigating to login');
      next({path: '/login'});
    } else {
      // Everything fine, go to page
      next();
    }
  }
});

export default router;

