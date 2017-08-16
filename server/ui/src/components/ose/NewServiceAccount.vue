<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">account_circle</i> Service-Account anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du für dein Projekt einen Service-Account anlegen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="createServiceAccount">
            <b-field label="Projekt-Name"
                     :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         name="Projekt-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <b-field label="Service-Account Name"
                     :type="errors.has('Service-Account') ? 'is-danger' : ''"
                     :message="errors.first('Service-Account')">
                <b-input v-model.trim="serviceAccount"
                         name="Service-Account"
                         v-validate="{ rules: { required: true, regex: /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/} }">
                </b-input>
            </b-field>
            <b-message type="is-info">
                Service-Account Name darf nur Kleinbuchstaben, Zahlen und - enthalten
            </b-message>

            <button :disabled="errors.any()"
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
      };
    },
    methods: {
      createServiceAccount: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
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
        });
      }
    }
  };
</script>