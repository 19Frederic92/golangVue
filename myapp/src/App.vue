<script setup>
// Aucun import de `Users.vue` ici pour laisser `RouterView` gérer l'affichage des pages.
</script>

<template>
  <div>
    <!-- Logo ou autre contenu global -->
    <a href="https://vitejs.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
    
    <!-- RouterView pour afficher le contenu des routes -->
    <RouterView />
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
