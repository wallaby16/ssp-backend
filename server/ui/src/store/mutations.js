export const STORAGE_KEY = 'cloud-ssp'

export const state = {
  user: JSON.parse(window.localStorage.getItem(STORAGE_KEY) || 'null'),
  notification: {}
}

export const mutations = {
  setUser(state, { user }) {
    state.user = user;
  },
  setNotification(state, { notification }) {
    state.notification = notification;
  }
}
