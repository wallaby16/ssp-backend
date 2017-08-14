<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Persistent Volume anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein Persistent Volume für OpenShift auf GlusterFS erstellen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="createGlusterVolume">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="projectname"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <p><em></em></p>
            <b-field label="Grösse">
                <b-input v-model.trim="size"
                         placeholder="100M"
                         required>
                </b-input>
            </b-field>
            <b-message type="is-warning">
                Grösse angeben mit Einheit (M/G) z.B. 100M oder 5G. Ab 1024M muss G verwendet werden
            </b-message>

            <b-field label="Name des Persistent Volume Claims">
                <b-input v-model.trim="pvcname"
                         required>
                </b-input>
            </b-field>

            <label class="label">Verwendungsmodus</label>
            <b-field>
                <b-radio-button v-model="mode"
                                native-value="ReadWriteOnce"
                                type="is-success">
                    <span>ReadWriteOnce (RWO)</span>
                </b-radio-button>

                <b-radio-button v-model="mode"
                                native-value="ReadWriteMany"
                                type="is-info">
                    <span>ReadWriteMany (RWX)</span>
                </b-radio-button>
            </b-field>
            <b-message type="is-warning">
                Siehe <a href="https://docs.openshift.com/container-platform/3.3/architecture/additional_concepts/storage.html#pv-access-modes">Dokumentation</a>
            </b-message>
            <br>

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Persistent Volume erstellen
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        projectname: '',
        pvcname: '',
        size: '',
        mode: 'ReadWriteOnce',
        loading: false,
        done: false
      }
    },
    methods: {
      createGlusterVolume: function(event) {
        this.loading = true;
//          this.done = true;
        // Todo do it
      }
    }
  }
</script>