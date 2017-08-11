export const STORAGE_KEY = 'cloud-ssp'

export const state = {
  user: JSON.parse(window.localStorage.getItem(STORAGE_KEY) || 'null')
}

export const mutations = {
  setUser(state, { user }) {
    state.user = user;
  }
}
