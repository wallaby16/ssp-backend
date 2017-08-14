import Vue from 'vue';
// Global Components
import App from './components/App.vue';
import Home from './components/shared/Home.vue';
import Nav from './components/shared/Nav.vue';
import Login from './components/shared/Login.vue';
import Notification from './components/shared/Notification.vue';
// OSE-Components
import EditQuota from './components/ose/EditQuota.vue';
import NewProject from './components/ose/NewProject.vue';
import NewTestProject from './components/ose/NewTestProject.vue';
import UpdateBilling from './components/ose/UpdateBilling.vue';
import NewServiceAccount from './components/ose/NewServiceAccount.vue';
// Gluster-Components
import FixVolume from './components/gluster/FixVolume.vue';
import NewVolume from './components/gluster/NewVolume.vue';
import GrowVolume from './components/gluster/GrowVolume.vue';

Vue.component('login', Login);
Vue.component('navbar', Nav);
Vue.component('notification', Notification);

export const GlobalComponents = {
  App,
  Login
};

export const LocalComponents = {
  Home,
  EditQuota,
  NewProject,
  NewTestProject,
  UpdateBilling,
  NewServiceAccount,
  FixVolume,
  NewVolume,
  GrowVolume
}



