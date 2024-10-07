<template>
  <v-form @submit.prevent>
    <v-row class="mb-3">
      <slot name="prepend-toolbar" />
      <v-spacer />
      <v-btn
        color="secondary"
        text="Reset"
        variant="text"
        prepend-icon="mdi-refresh"
        @click="emit('reset')"
      />
      <v-btn
        color="primary"
        text="Save settings"
        variant="text"
        prepend-icon="mdi-floppy"
        @click="emit('submit')"
      />
    </v-row>
    <v-divider class="mb-3" />
    <slot />
  </v-form>
</template>

<script setup lang="ts" generic="Settings extends {}, P extends SchemaProperties">
import { useSlots } from 'vue'
import { SchemaProperties } from '../toolkit/forms/form'

useSlots()

const emit = defineEmits<{
  reset: []
  submit: []
}>()
</script>

<style scoped></style>
