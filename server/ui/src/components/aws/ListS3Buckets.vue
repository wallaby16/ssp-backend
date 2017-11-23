<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">cached</i> AWS S3 Bucket anzeigen</h1>
                </div>
                <h2 class="subtitle">
                    Hier werden alle deine AWS S3 Buckets angezeigt.</h2>
            </div>
        </div>
        <br>
        <b-table :data="data"
                 v-bind:class="{'is-loading': loading}"
                 :narrowed="true">

            <template scope="props">
                <b-table-column field="name" label="Bucket-Name">
                    {{ props.row.name }}
                </b-table-column>
                <b-table-column field="account" label="SBB AWS Account">
                    {{ props.row.account }}
                </b-table-column>
            </template>

            <div slot="empty" class="has-text-centered">
                Hier werden deine Buckets angezeigt wenn du welche hast
            </div>

        </b-table>
    </div>
</template>

<script>
  export default {
    data() {
      return {
        data: [],
        loading: false
      };
    },
    mounted: function() {
        this.listS3Buckets();
    },
    methods: {
      listS3Buckets: function() {
        this.loading = true;
        this.$http.get('/api/aws/s3').then((res) => {
          this.data = res.body.buckets;
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    }
  };
</script>