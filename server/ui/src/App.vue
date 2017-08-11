<template>
    <div>
        <navbar></navbar>
        <div class="container">
            <router-view></router-view>
        </div>
    </div>
</template>

<style lang="scss">
    @import "~bulma/sass/utilities/_all";

    $navbar-background-color: hsl(217, 71%, 53%);
    $navbar-item-color: white;
    $navbar-item-active-color: $navbar-item-color;
    $navbar-dropdown-background-color: $navbar-background-color;

    // Import Bulma and Buefy styles
    @import "~bulma";
    @import "~buefy/src/scss/buefy";
</style>

