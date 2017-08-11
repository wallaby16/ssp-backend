<template>
    <section class="section">
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Projekt Quotas anpassen</h1>
                </div>
                <h2 class="subtitle">Hier kannst du CPU / Memory Quotas deines Projektes anpassen. Alle Anpassungen werden geloggt.</h2>
            </div>
        </div>
        <div class="container">
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

            <button class="button is-primary" v-on:click.once="createProject">Projekt erstellen</button>
        </div>
    </section>
</template>

<script>
    module.exports = {
      data: function() {
        return {
          cpu: 2,
          memory: 4,
          projectname: ''
        }
      },
      methods: {
        createProject: function(event) {
          if (this.projectname && this.cpu && this.memory) {
            // Todo do it
          } else {
            this.$toast.open({
              type: 'is-danger',
              message: "Bitte alle Felder ausfüllen"
            })
          }
        }
      }
    }
</script>