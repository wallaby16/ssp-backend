<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">
                        <i class="material-icons">receipt</i> OpenShift Test-Projekt anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein OpenShift Test-Projekt erstellen</h2>
            </div>
        </div>
        <br>
        <article class="message is-info">
            <div class="message-body">
                - Du kannst jederzeit ein Testprojekt erstellen um etwas auszuprobieren. Es entstehen keine Kosten.<br/>
                - Ein Test-Projekt enthält deine u-Nummer<br/>
                - Pods aus einem Test-Projekt können jederzeit durch das Cloud-Team gestoppt werden<br/>
                - Bitte das Testprojekt nach Ende des Tests selbst löschen<br/>
                - Nicht gelöschte Testprojekte werden durch das Cloud-Team gelöscht<br/>
            </div>
        </article>
        <form v-on:submit.prevent="newTestProject">
            <b-field>
                <label class="label">Testprojekt-Name</label>
            </b-field>
            <b-field class="has-addons" :type="errors.has('Testprojekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Testprojekt-Name')">
                <p class="control">
                    <span class="button is-static">{{ username }}-</span>
                </p>
                <p class="control">
                    <b-input v-model.trim="testprojectname" name="Testprojekt-Name"
                             v-validate="{ rules: { required: true, regex: /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/} }"
                             placeholder="testprojekt">
                    </b-input>
                </p>
            </b-field>

            <b-message type="is-info">
                Testprojekt-Name darf nur Kleinbuchstaben, Zahlen und - enthalten
            </b-message>

            <button :disabled="errors.any()" v-bind:class="{'is-loading': loading}" class="button is-primary">
                Neues Test-Projekt erstellen
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
      };
    },
    methods: {
      newTestProject: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/ose/testproject', {
              project: this.testprojectname
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