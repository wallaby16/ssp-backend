<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">perm_scan_wifi</i> Gluster Konfiguration erzeugen</h1>
                </div>
                <h2 class="subtitle">
                    Diese Funktion erstellt die Gluster Objekte (Service & Endpoints) in deinem Projekt</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="fixGlusterObjects">
            <b-field label="Projekt-Name"
                :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         name="Projekt-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <button :disabled="errors.any()"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Gluster Objekte erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        project: '',
        loading: false
      };
    },
    methods: {
      fixGlusterObjects: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/gluster/volume/fix', {
              project: this.project
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
