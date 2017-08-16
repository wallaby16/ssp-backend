<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">speaker_notes</i> OpenShift Projekt anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein OpenShift Projekt erstellen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="newProject">
            <b-field label="Projekt-Name"
                     :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         name="Projekt-Name"
                         v-validate="{ rules: { required: true, regex: /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/} }"
                         placeholder="projekt-dev">
                </b-input>
            </b-field>
            <b-message type="is-info">
                Projekt-Name darf nur Kleinbuchstaben, Zahlen und - enthalten
            </b-message>

            <b-field label="Kontierungsnummer"
                     :type="errors.has('Kontierungsnummer') ? 'is-danger' : ''"
                     :message="errors.first('Kontierungsnummer')">
                <b-input v-model.trim="billing"
                         name="Kontierungsnummer"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <b-field label="MEGA-ID">
                <b-input v-model.trim="megaId"></b-input>
            </b-field>

            <button :disabled="errors.any()"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Neues Projekt erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        megaId: '',
        billing: '',
        project: '',
        loading: false
      };
    },
    methods: {
      newProject: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/ose/project', {
              project: this.project,
              billing: this.billing,
              megaId: this.megaId
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