<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Projekt Quotas anpassen</h1>
                </div>
                <h2 class="subtitle">
                    Hier kannst du CPU / Memory Quotas deines Projektes anpassen. Alle Anpassungen werden geloggt.</h2>
            </div>
        </div>
        <br>
        <form v-on:submit.prevent="editQuotas">
            <b-field label="Projekt-Name">
                <b-input v-model.trim="projectname" required></b-input>
            </b-field>

            <b-field label="Neue CPU Quotas [Cores]">
                <b-input type="number"
                         v-model.number="cpu"
                         min="1">
                </b-input>
            </b-field>

            <b-field label="Neue Memory Quotas [GB]">
                <b-input type="number"
                         v-model.number="memory"
                         min="2">
                </b-input>
            </b-field>

            <button type="submit" v-if="!done"
                    v-bind:class="{'is-loading': loading}"
                    class="button is-primary">Quotas anpassen
            </button>
        </form>
    </section>
</template>

<script>
  export default {
    data() {
      return {
        cpu: 2,
        memory: 4,
        projectname: '',
        loading: false,
        done: false
      }
    },
    methods: {
      editQuotas: function(event) {
        this.loading = true;
//          this.done = true;
        // Todo do it
      }
    }
  }
</script>