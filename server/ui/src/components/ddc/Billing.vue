<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><i class="material-icons">attach_money</i> DDC-Kosten anzeigen</h1>
                </div>
                <h2 class="subtitle">
                    Hier findest du die Kosten für deine DDC-Infrastrukturen</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="getDDCBilling">
            <button type="submit"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Kostenberechnung erstellen
            </button>
        </form>
        <br>
        <b-table :data="data"
                 :narrowed="true">

            <template scope="props">
                <b-table-column field="sender" label="Von" width="40">
                    {{ props.row.sender }}
                </b-table-column>
                <b-table-column field="assignment" label="Nach" width="40">
                    {{ props.row.assignment }}
                </b-table-column>
                <b-table-column field="art" label="Art" width="40">
                    {{ props.row.art }}
                </b-table-column>
                <b-table-column field="project" label="Projekt" width="40">
                    {{ props.row.project }}
                </b-table-column>
                <b-table-column field="host" label="Host" width="40">
                    {{ props.row.host }}
                </b-table-column>
                <b-table-column field="total" label="Total" width="40" numeric>
                    {{ props.row.total }} CHF
                </b-table-column>
            </template>

            <div slot="empty" class="has-text-centered">
                Klicke auf den Button oben um die Kosten zu berechnen
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
      }
    },
    methods: {
      getDDCBilling: function() {
        this.loading = true;

        this.$http.get('/api/ddc/billing').then((res) => {
          this.data = res.body;
          this.loading = false;
        }, () => {
          this.loading = false;
        });
      }
    }
  }
</script>