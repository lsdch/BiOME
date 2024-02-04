<template>
  <v-snackbar v-model="active" :timeout="timeout" :color="color">
    {{ text }}
    <template v-slot:actions>
      <v-btn variant="text" @click="active = false"> Close </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
const active = defineModel<boolean>()

export type FeedbackProps = {
  color?: string
  text: string
  timeout?: number
}

withDefaults(defineProps<FeedbackProps>(), { timeout: 3000 })
</script>

<style scoped></style>
