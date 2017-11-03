<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">
                        <i class="material-icons">person</i> OpenShift Projekt Admins anzeigen</h1>
                    <h2 class="subtitle">
                        Zeigt alle User eines Projektes mit Admin-Rolle an</h2>
                </div>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="getAdminList">
            <b-field>
                <label class="label">Projekt-Name</label>
            </b-field>
            <b-field :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="projectname" name="Projekt-Name"
                         v-validate="{ rules: { required: true, regex: /^[a-z0-9]([-a-z0-9]*[a-z0-9])$/} }"
                         placeholder="projekt">
                </b-input>
            </b-field>

            <button :disabled="errors.any()" v-bind:class="{'is-loading': loading}" class="button is-primary">
                Admin-Liste anzeigen
            </button>
        </form>
        <br><br>
        <label><strong>Administratoren</strong></label><br><br>
        <b-table :data="data"
                 :narrowed="true">

            <template scope="props">
                <b-table-column field="sender" label="Benutzername">
                    {{ props.row }}
                </b-table-column>
            </template>

            <div slot="empty" class="has-text-centered">
                Noch keine Abfrage durchgeführt
            </div>
        </b-table>
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
        data: [],
        loading: false
      };
    },
    methods: {
      getAdminList: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.get('/api/ose/project/' + this.projectname + '/admins').then((res) => {
              this.loading = false;

              this.data = res.body.admins;
            }, () => {
              this.loading = false;
            });
          }
        });
      }
    }
  };
</script>