<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">perm_identity</i> AWS S3 Bucket-Benutzer erstellen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du einen Benutzer zu einem AWS S3 Bucket erstellen. Alle Bestellungen werden geloggt & verrechnet.</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="newS3User">
            <b-field label="Bucket-Name"
                     :type="errors.has('Bucket-Name') ? 'is-danger' : ''"
                     :message="errors.first('Bucket-Name')">

                <b-select placeholder="Wähle den Bucket aus"
                          :loading="loading"
                          v-model="bucket"
                          required>
                    <option
                            v-for="bucket in buckets"
                            :value="bucket.name"
                            :key="bucket.name">
                        {{ bucket.name }}
                    </option>
                </b-select>
            </b-field>

            <b-field label="Name vom neuen Benutzer"
                     :type="errors.has('Name vom neuen Benutzer') ? 'is-danger' : ''"
                     :message="errors.first('Name vom neuen Benutzer')">
                <b-input type="text"
                         v-validate="{ rules: { required: true, regex: /^[a-zA-Z0-9\-]+$/} }"
                         name="Name vom neuen Benutzer"
                         v-model.number="username">
                </b-input>
            </b-field>

            <label class="label">Rechte: Lesen / Schreiben</label>
            <b-field>
                <b-radio-button v-model="isReadonly"
                                native-value="true"
                                type="is-success">
                    <span>Nur lesen</span>
                </b-radio-button>
                <b-radio-button v-model="isReadonly"
                                native-value="false"
                                type="is-danger">
                    <span>Lesen & Schreiben</span>
                </b-radio-button>
            </b-field>

            <button :disabled="errors.any() || loading"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Benutzer zu S3 Bucket erstellen
            </button>
        </form>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        buckets: [],
        username: '',
        bucket: '',
        isReadonly: "true",
        loading: false
      };
    },
    mounted: function() {
      this.getUsersBuckets();
    },
    methods: {
      getUsersBuckets: function() {
        this.loading = true;
        this.$http.get('/api/aws/s3').then((res) => {
          this.buckets = res.body.buckets;
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      },
      newS3User: function() {
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.loading = true;

            this.$http.post('/api/aws/s3/' + this.bucket + '/user', {
              username: this.username,
              isReadonly: this.isReadonly == "true"
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