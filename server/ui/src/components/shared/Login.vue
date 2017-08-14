<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Login</h1>
                </div>
            </div>
        </div>
        <form v-on:submit.prevent="login">
            <b-field label="Benutzername (U-Nummer)">
                <b-input v-model.trim="username" required></b-input>
            </b-field>

            <b-field label="Passwort">
                <b-input type="password" v-model="password" required></b-input>
            </b-field>

            <button type="submit"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Login
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        username: '',
        password: '',
        loading: false
      }
    },
    methods: {
      login: function() {
        this.loading = true;

        this.$http.post('/login', {
          username: this.username,
          password: this.password
        }).then(res => {
          this.loading = false;

          // Decode JWT
          const base64Url = res.body.token.split('.')[1];
          const base64 = base64Url.replace('-', '+').replace('_', '/');
          const userData = JSON.parse(window.atob(base64));

          this.$store.commit('setUser', {
            user: {
              name: userData.id,
              token: res.body.token
            }
          });

          this.$toast.open({
            type: 'is-success',
            message: 'Login war erfolgreich'
          })

          this.$router.push({path: '/'})

        }, () => {
          this.loading = false;
        });
      }
    }
  }
</script>