<template>
  <v-icon
    v-bind="
      exists
        ? {
            icon: 'mdi-link-variant',
            color: 'primary',
            title: 'Existing site'
          }
        : {
            icon: 'mdi-plus-thick',
            color: 'warning',
            title: 'New site'
          }
    "
  />
</template>

<script setup lang="ts">
defineProps<{ exists?: boolean }>()
</script>

<style scoped></style>
