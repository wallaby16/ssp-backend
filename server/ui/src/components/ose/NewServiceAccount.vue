<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Service-Account anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du für dein Projekt einen Service-Account anlegen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="createServiceAccount">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <b-field label="Service-Account Name">
                <b-input v-model.trim="serviceAccount"
                         required>
                </b-input>
            </b-field>

            <button type="submit"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Service-Account erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        serviceAccount: '',
        project: '',
        loading: false
      }
    },
    methods: {
      createServiceAccount: function() {
        this.loading = true;

        this.$http.post('/api/ose/serviceaccount', {
          project: this.project,
          serviceAccount: this.serviceAccount
        }).then(() => {
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    }
  }
</script>