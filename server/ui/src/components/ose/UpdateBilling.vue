<template>
    <div>
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
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <b-field label="Kontierungsnummer">
                <b-input v-model.trim="billing"
                         required>
                </b-input>
            </b-field>

            <button type="submit"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Kontierungsinformation anpassen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        billing: '',
        project: '',
        loading: false
      }
    },
    methods: {
      updateBilling: function() {
        this.loading = true;

        this.$http.post('/api/ose/billing', {
          project: this.project,
          billing: this.billing
        }).then(() => {
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    }
  }
</script>