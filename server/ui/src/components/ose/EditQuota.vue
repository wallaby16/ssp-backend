<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">cached</i> Projekt Quotas anpassen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du CPU / Memory Quotas deines Projektes anpassen. Alle Anpassungen werden geloggt.</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="editQuotas">
            <b-field label="Projekt-Name"
                     :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         name="Projekt-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <b-field label="Neue CPU Quotas [Cores]"
                     :type="errors.has('CPU') ? 'is-danger' : ''"
                     :message="errors.first('CPU')">
                <b-input type="number"
                         v-validate="'required'"
                         name="CPU"
                         v-model.number="cpu"
                         min="1">
                </b-input>
            </b-field>

            <b-field label="Neue Memory Quotas [GB]"
                     :type="errors.has('Memory') ? 'is-danger' : ''"
                     :message="errors.first('Memory')">
                <b-input type="number"
                         v-model.number="memory"
                         v-validate="'required'"
                         name="Memory"
                         min="2">
                </b-input>
            </b-field>

            <button :disabled="errors.any()"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Quotas anpassen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        cpu: 2,
        memory: 4,
        project: '',
        loading: false
      };
    },
    methods: {
      editQuotas: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/ose/quotas', {
              project: this.project,
              cpu: '' + this.cpu,
              memory: '' + this.memory
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