<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">cached</i> AWS S3 Bucket erstellen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du einen AWS S3 Bucket erstellen. Alle Bestellungen werden geloggt & verrechnet.</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="newS3Bucket">
            <b-field label="Projekt-Name"
                     :type="errors.has('Projekt-Name') ? 'is-danger' : ''"
                     :message="errors.first('Projekt-Name')">
                <b-input v-model.trim="project"
                         name="Projekt-Name"
                         v-validate="'required'">
                </b-input>
            </b-field>

            <b-field label="Bucket-Name"
                     :type="errors.has('Bucket-Name') ? 'is-danger' : ''"
                     :message="errors.first('Bucket-Name')">
                <b-input type="text"
                         v-validate="{ rules: { required: true, regex: /^[a-zA-Z0-9\-]+$/} }"
                         name="Bucket-Name"
                         v-model.number="bucketname">
                </b-input>
            </b-field>
            <b-message type="is-info">
                Bucket Name wird um folgendes ergänzt: sbb-"dein name"-stage
            </b-message>

            <b-field label="Kontierungsnummer"
                     :type="errors.has('Kontierungsnummer') ? 'is-danger' : ''"
                     :message="errors.first('Kontierungsnummer')">
                <b-input type="text"
                         v-model.number="billing"
                         v-validate="'required'"
                         name="Kontierungsnummer">
                </b-input>
            </b-field>

            <label class="label">SBB AWS Account</label>
            <b-field>
                <b-radio-button v-model="stage"
                                native-value="dev"
                                type="is-success">
                    <span>Entwicklung</span>
                </b-radio-button>
                <b-radio-button v-model="stage"
                                native-value="test"
                                type="is-success">
                    <span>Test</span>
                </b-radio-button>
                <b-radio-button v-model="stage"
                                native-value="int"
                                type="is-success">
                    <span>Integration</span>
                </b-radio-button>
                <b-radio-button v-model="stage"
                                native-value="prod"
                                type="is-info">
                    <span>Produktion</span>
                </b-radio-button>
            </b-field>

            <button :disabled="errors.any()"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">S3 Bucket erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        project: '',
        bucketname: '',
        billing: '',
        stage: 'dev',
        loading: false
      };
    },
    methods: {
      newS3Bucket: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/aws/s3', {
              project: this.project,
              bucketname: this.bucketname,
              billing: '' + this.billing,
              stage: '' + this.stage
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