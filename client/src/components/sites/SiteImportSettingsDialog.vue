<template>
  <FormDialog v-model="openState" title="Import settings">
    <CoordPrecisionPicker label="Default coordinates precision" />
  </FormDialog>
  <slot name="activator" :open :close></slot>
</template>

<script setup lang="ts">
import FormDialog from '../toolkit/forms/FormDialog.vue'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'

const openState = defineModel<boolean>('open', { default: false })

function open() {
  openState.value = true
}

function close() {
  openState.value = false
}
</script>

<style scoped></style>
