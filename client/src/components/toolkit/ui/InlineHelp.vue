<template>
  <v-menu location="right center" origin="start center" offset="20">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        variant="plain"
        class="d-flex align-center justify-center"
        icon="mdi-help-circle"
        density="compact"
        :rounded="100"
        color="info"
        size="small"
      >
      </v-btn>
    </template>
    <v-card :max-width="300">
      <v-card-text>
        <slot />
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts"></script>

<style scoped lang="scss"></style>
