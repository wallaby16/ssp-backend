<template>
    <section class="section">
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
                <p><em>Projekt-Name darf nur Kleinbuchstaben, Zahlen und - enthalten</em></p><br>
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

                <b-field label="MEGA-ID">
                    <b-input v-model.trim="megaid"
                             required>
                    </b-input>
                </b-field>

                <button type="submit" v-if="!done"
                        v-bind:class="{'is-loading': loading}"
                        class="button is-primary">Neues Projekt erstellen</button>
            </form>
    </section>
</template>

<script>
  module.exports = {
    data: function() {
      return {
        megaid: '',
        billing: '',
        projectname: '',
        loading: false,
        done: false
      }
    },
    methods: {
      newProject: function(event) {
        this.projectname = this.projectname.toLowerCase()
        if (this.projectname && this.megaid && this.billing) {
          this.loading = true;
//          this.done = true;
          // Todo do it
        } else {
          this.$toast.open({
            type: 'is-danger',
            message: "Bitte alle Felder ausfüllen"
          })
        }
      }
    }
  }
</script>