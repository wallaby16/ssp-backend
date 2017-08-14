import {LocalComponents} from './components';
import VueRouter from 'vue-router';

const routes = [
  {
    path: '/', component: LocalComponents.Home
  },
  {
    path: '/login', component: LocalComponents.Login
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
  }
];

export default new VueRouter({routes});

