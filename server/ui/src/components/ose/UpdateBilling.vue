<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Kontierungsnummer anpassen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du die Kontierungsnummer deines OpenShift Projekts anpassen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="updateBilling">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="projectname"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <b-field label="Kontierungsnummer">
                <b-input v-model.trim="billing"
                         required>
                </b-input>
            </b-field>

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Kontierungsinformation anpassen
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        billing: '',
        projectname: '',
        loading: false,
        done: false
      }
    },
    methods: {
      updateBilling: function(event) {
        this.loading = true;
//          this.done = true;
        // Todo do it
      }
    }
  }
</script>