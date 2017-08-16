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
        <b-message type="is-danger">
            In OpenShift wird nach dem Vergrössern immer noch die alte Grösse angegeben sein. Dieser Wert lässt sich im Moment leider nicht verändern.
        </b-message>

        <form v-on:submit.prevent="growGlusterVolume">
            <b-field label="Projekt-Name"
                     :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         name="Projekt-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <b-field label="Neue Grösse"
                     :type="errors.has('Grösse') ? 'is-danger' : ''"
                     :message="errors.first('Grösse')">
                <b-input v-model.trim="newSize"
                         placeholder="100M"
                         name="Grösse"
                         v-validate="'required'">
                </b-input>
            </b-field>
            <b-message type="is-info">
                Das Volume wird auf die angegebene Grösse vergrösert. Verkleinern ist nicht möglich. z.B. 100M oder 5G
            </b-message>

            <p><em></em></p>
            <b-field label="Name des Persistent Volumes"
                     :type="errors.has('PV-Name') ? 'is-danger' : ''"
                     :message="errors.first('PV-Name')">
                <b-input v-model.trim="pvName"
                         name="PV-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>
            <b-message type="is-info">
                Nicht der Name des PVC, sondern das was in OpenShift unter "Storage" > Spalte "Status" > <strong>fett</strong>
                geschrieben ist
            </b-message>

            <button :disabled="errors.any()"
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
      };
    },
    methods: {
      growGlusterVolume: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
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
        });
      }
    }
  };
</script>