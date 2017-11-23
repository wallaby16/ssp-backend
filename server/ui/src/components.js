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
import AdminList from './components/ose/AdminList.vue';
import UpdateBilling from './components/ose/UpdateBilling.vue';
import NewServiceAccount from './components/ose/NewServiceAccount.vue';
// Gluster-Components
import FixVolume from './components/gluster/FixVolume.vue';
import NewVolume from './components/gluster/NewVolume.vue';
import GrowVolume from './components/gluster/GrowVolume.vue';
// DDC-Components
import DDCBilling from './components/ddc/Billing.vue';
// AWS Components
import ListS3Buckets from './components/aws/ListS3Buckets.vue';
import NewS3Bucket from './components/aws/NewS3Bucket.vue';

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
  AdminList,
  NewTestProject,
  UpdateBilling,
  NewServiceAccount,
  FixVolume,
  NewVolume,
  GrowVolume,
  DDCBilling,
  ListS3Buckets,
  NewS3Bucket
}



