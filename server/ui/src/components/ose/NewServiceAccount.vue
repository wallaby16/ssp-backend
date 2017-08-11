<template>
    <section class="section">
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
                <b-input v-model.trim="projectname"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <b-field label="Service-Account Name">
                <b-input v-model.trim="serviceaccount"
                         required>
                </b-input>
            </b-field>

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Service-Account erstellen
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        serviceaccount: '',
        projectname: '',
        loading: false,
        done: false
      }
    },
    methods: {
      createServiceAccount: function(event) {
        this.loading = true;
//          this.done = true;
        // Todo do it
      }
    }
  }
</script>