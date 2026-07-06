<template>
  <v-card :title :subtitle class="small-card-title border-md" variant="outlined" flat>
    <template #append>
      <slot name="append">
        <v-switch v-model="active" color="primary" hide-details density="compact"></v-switch>
      </slot>
    </template>
    <template #default v-if="active">
      <v-divider></v-divider>
      <slot name="default"></slot>
    </template>
  </v-card>
</template>

<script setup lang="ts">
const active = defineModel<boolean>({ default: false })
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<style scoped lang="scss"></style>
