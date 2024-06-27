<template>
  <FormDialog v-model="open" title="Import settings">
    <CoordPrecisionPicker label="Default coordinates precision" />
  </FormDialog>
</template>

<script setup lang="ts">
import FormDialog from '../toolkit/forms/FormDialog.vue'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'

const open = defineModel<boolean>()
</script>

<style scoped></style>
