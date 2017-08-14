<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> OpenShift Projekt anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein OpenShift Projekt erstellen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="newProject">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>
            <b-message type="is-info">
                Projekt-Name darf nur Kleinbuchstaben, Zahlen und - enthalten
            </b-message>

            <b-field label="Kontierungsnummer">
                <b-input v-model.trim="billing"
                         required>
                </b-input>
            </b-field>

            <b-field label="MEGA-ID">
                <b-input v-model.trim="megaId"
                         required>
                </b-input>
            </b-field>

            <button type="submit"
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
      }
    },
    methods: {
      newProject: function() {
        this.project = this.project.toLowerCase()
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
    }
  }
</script>