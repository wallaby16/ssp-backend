<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">cloud_upload</i> Persistent Volume vergrösern</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein Persistent Volume vergrössern</h2>
            </div>
        </div>
        <br>
        <b-message title="ACHTUNG" type="is-danger">
            In OpenShift wird nach dem Vergrössern immer noch die alte Grösse angegeben sein. Dieser Wert lässt sich im Moment leider nicht verändern.
        </b-message>

        <form v-on:submit.prevent="growGlusterVolume">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         required>
                </b-input>
            </b-field>

            <b-field label="Neue Grösse">
                <b-input v-model.trim="newSize"
                         placeholder="100M"
                         required>
                </b-input>
            </b-field>
            <b-message type="is-info">
                Das Volume wird auf die angegebene Grösse vergrösert. Verkleinern ist nicht möglich. z.B. 100M oder 5G
            </b-message>

            <p><em></em></p>
            <b-field label="Name des Persistent Volumes">
                <b-input v-model.trim="pvName"
                         required>
                </b-input>
            </b-field>
            <b-message type="is-info">
                Nicht der Name des PVC, sondern das was in OpenShift unter "Storage" > Spalte "Status" > <strong>fett</strong> geschrieben ist
            </b-message>

            <button type="submit"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Persistent Volume vergrössern
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        project: '',
        pvName: '',
        newSize: '',
        loading: false
      }
    },
    methods: {
      growGlusterVolume: function() {
        this.loading = true;

        this.$http.post('/api/gluster/volume/grow', {
          project: this.project,
          newSize: this.newSize,
          pvName: this.pvName
        }).then(() => {
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    }
  }
</script>