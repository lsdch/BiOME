<template>
  <v-snackbar v-model="open" vertical variant="flat" :color="colors.red.darken1">
    <div class="text-overline mb-3 d-flex align-center">
      <v-icon class="mr-2" icon="mdi-alert-circle" color="white" size="large" />
      {{ title }}
    </div>
    <div v-for="err in errors" :key="err.message">
      {{ err.message }}
    </div>
    <template v-slot:actions>
      <v-btn color="white" variant="text" @click="open = false"> Close </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ErrorDetail } from '@/api'
import colors from 'vuetify/util/colors'

const open = defineModel<boolean>({ required: true })

defineProps<{
  title: string
  errors?: ErrorDetail[]
}>()
</script>

<style scoped></style>
