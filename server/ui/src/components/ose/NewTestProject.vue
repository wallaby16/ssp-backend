<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> OpenShift Test-Projekt anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein OpenShift Test-Projekt erstellen</h2>
            </div>
        </div>
        <br>
            <form v-on:submit.prevent="newTestProject">
                <p><em>Projekt-Name darf nur Kleinbuchstaben, Zahlen und - enthalten</em></p><br>


                <b-field>
                    <label class="label">Projekt-Name</label>
                </b-field>

                <b-field class="has-addons">
                    <p class="control">
                        <span class="button is-static">{{ userId }}-</span>
                    </p>
                    <p class="control">
                        <b-input v-model.trim="testprojectname"
                                 placeholder="testprojekt"
                                 required>
                        </b-input>
                    </p>
                </b-field>

                <button type="submit" v-if="!done"
                        v-bind:class="{'is-loading': loading}"
                        class="button is-primary">Neues Test-Projekt erstellen</button>
            </form>
    </section>
</template>

<script>
  module.exports = {
    data: function() {
      return {
        userId: 'u220374',
        testprojectname: '',
        loading: false,
        done: false
      }
    },
    methods: {
      newTestProject: function(event) {
        this.testprojectname = this.testprojectname.toLowerCase()
        if (this.testprojectname) {
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