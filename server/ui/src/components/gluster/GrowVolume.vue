<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Persistent Volume vergrösern</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein Persistent Volume vergrössern</h2>
            </div>
        </div>
        <br>
        <b-message type="is-danger">
            ACHTUNG: In OpenShift wird nach dem Vergrössern immer noch die alte Grösse angegeben sein. Dieser Wert lässt sich im Moment leider nicht verändern.
        </b-message>

        <form v-on:submit.prevent="growGlusterVolume">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="projectname"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <b-field label="Neue Grösse">
                <b-input v-model.trim="newsize"
                         placeholder="100M"
                         required>
                </b-input>
            </b-field>
            <b-message type="is-warning">
                Das Volume wird auf die angegebene Grösse vergrösert. Verkleinern ist nicht möglich. z.B. 100M oder 5G
            </b-message>

            <p><em></em></p>
            <b-field label="Name des Persistent Volumes">
                <b-input v-model.trim="pvname"
                         required>
                </b-input>
            </b-field>
            <b-message type="is-warning">
                Nicht der Name des PVC, sondern das was in OpenShift unter "Storage" > Spalte "Status" > <strong>fett</strong> geschrieben ist
            </b-message>

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Persistent Volume vergrössern
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        projectname: '',
        pvname: '',
        newsize: '',
        loading: false,
        done: false
      }
    },
    methods: {
      growGlusterVolume: function(event) {
        this.loading = true;
//          this.done = true;
        // Todo do it
      }
    }
  }
</script>