<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Gluster Konfiguration erzeugen</h1>
                </div>
                <h2 class="subtitle">
                    Diese Funktion erstellt die Gluster Objekte (Service & Endpoints) in deinem Projekt</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="fixGlusterObjects">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="projectname"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Gluster Objekte erstellen
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        projectname: '',
        loading: false,
        done: false
      }
    },
    methods: {
      fixGlusterObjects: function(event) {
        this.loading = true;
//          this.done = true;
        // Todo do it
      }
    }
  }
</script>
