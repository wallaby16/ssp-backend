<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">cloud</i> Persistent Volume anlegen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du ein Persistent Volume für OpenShift auf GlusterFS erstellen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="createGlusterVolume">
            <b-field label="Projekt-Name"
                     :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         placeholder="projekt-dev"
                         name="Projekt-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <b-field label="Grösse"
                     :type="errors.has('Grösse') ? 'is-danger' : ''"
                     :message="errors.first('Grösse')">
                <b-input v-model.trim="size"
                         placeholder="100M"
                         name="Grösse"
                         v-validate="'required'">
                </b-input>
            </b-field>
            <b-message type="is-info">
                Grösse angeben mit Einheit (M/G) z.B. 100M oder 5G. Ab 1024M muss G verwendet werden
            </b-message>

            <b-field label="Name des Persistent Volume Claims"
                     :type="errors.has('PVC-Name') ? 'is-danger' : ''"
                     :message="errors.first('PVC-Name')">
                <b-input v-model.trim="pvcName"
                         name="PVC-Name"
                         v-validate="'required'">
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
            <b-message type="is-info">
                Siehe <a
                    href="https://docs.openshift.com/container-platform/3.3/architecture/additional_concepts/storage.html#pv-access-modes">Dokumentation</a>
            </b-message>
            <br>

            <button :disabled="errors.any()"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Persistent Volume erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        project: '',
        pvcName: '',
        size: '',
        mode: 'ReadWriteOnce',
        loading: false
      };
    },
    methods: {
      createGlusterVolume: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/gluster/volume', {
              project: this.project,
              size: this.size,
              pvcName: this.pvcName,
              mode: this.mode
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