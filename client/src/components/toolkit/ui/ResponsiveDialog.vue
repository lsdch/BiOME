<template>
  <v-dialog v-model="dialogOpen" fullscreen :transition="false">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>

    <template #default="{ isActive }">
      <slot v-if="asDialog && isActive.value" name="default" :is-dialog="asDialog" :is-active />
    </template>
  </v-dialog>
  <slot v-if="!asDialog" name="default" :is-dialog="asDialog" :is-active="false" />
</template>

<script setup lang="ts">
defineProps<{ asDialog: boolean }>()

const dialogOpen = defineModel<boolean>('open')
</script>

<style scoped></style>
