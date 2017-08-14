<template>
    <section class="section">
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

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Login
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        username: '',
        password: '',
        loading: false,
        done: false
      }
    },
    methods: {
      login: function() {
        this.loading = true;
        this.done = true;

        this.$store.commit('setUser', {user: {name: 'u220374'}});

        this.$toast.open({
          type: 'is-success',
          message: 'Login war erfolgreich'
        })

        this.$router.push({path: '/'})
      }
    }
  }
</script>