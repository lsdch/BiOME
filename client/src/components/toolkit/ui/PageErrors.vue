<template>
  <v-container>
    <v-alert color="error" v-bind="errorsMap[error.status ?? 500]" />
  </v-container>
</template>

<script setup lang="ts" generic="T">
import { ErrorModel } from '@/api'

defineProps<{ error: ErrorModel }>()

type ErrorInfos = {
  icon: string
  text: string
}
const errorsMap: Record<number, ErrorInfos> = {
  404: { icon: 'mdi-close-circle', text: 'Item does not exist.' },
  401: { icon: 'mdi-lock', text: 'Access denied.' },
  500: { icon: 'mdi-alert-box', text: 'Unexpected server error.' }
}
</script>

<style scoped lang="scss"></style>
