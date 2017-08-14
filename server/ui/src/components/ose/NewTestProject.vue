<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">receipt</i> OpenShift Test-Projekt anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein OpenShift Test-Projekt erstellen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="newTestProject">
            <b-field>
                <label class="label">Testprojekt-Name</label>
            </b-field>
            <b-field class="has-addons">
                <p class="control">
                    <span class="button is-static">{{ username }}-</span>
                </p>
                <p class="control">
                    <b-input v-model.trim="testprojectname"
                             placeholder="testprojekt"
                             required>
                    </b-input>
                </p>
            </b-field>

            <b-message type="is-info">
                Testprojekt-Name darf nur Kleinbuchstaben, Zahlen und - enthalten
            </b-message>

            <button type="submit"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Neues Test-Projekt erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    computed: {
      username() {
        return this.$store.state.user.name;
      }
    },
    data() {
      return {
        testprojectname: '',
        loading: false
      }
    },
    methods: {
      newTestProject: function() {
        this.testprojectname = this.testprojectname.toLowerCase()
        this.loading = true;

        this.$http.post('/api/ose/testproject', {
          project: this.testprojectname
        }).then(() => {
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    }
  }
</script>