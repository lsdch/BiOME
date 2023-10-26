<template>
  <div class="w-100 d-flex align-center flex-column mb-5">
    <v-btn
      class="plain-link text-h3 font-weight-bold"
      :to="{ name: 'home' }"
      :text="app_name"
      variant="text"
      size="x-large"
      :ripple="false"
      plain
    />
  </div>
</template>

<script setup lang="ts">
const app_name = import.meta.env.VITE_APP_NAME
</script>

<style>
.plain-link .v-btn__overlay {
  display: none;
}
</style>
