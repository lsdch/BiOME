<template>
  <v-dialog v-model="isRevealed" max-width="500px" persistent @keyup.esc="emit('cancel')">
    <v-card>
      <v-toolbar dark dense flat>
        <v-toolbar-title class="text-body-2 font-weight-bold grey--text">
          {{ title }}
        </v-toolbar-title>
      </v-toolbar>
      <v-card-text v-if="message"> {{ message }} </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="grey" variant="text" @click="emit('cancel')" text="Cancel" />
        <v-btn color="blue-darken-1" variant="text" @click="emit('confirm')" text="OK" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="Payload = any">
export type ConfirmDialogProps<T> = {
  title: string
  message?: string
  data?: T
}

const isRevealed = defineModel<boolean>()

defineProps<ConfirmDialogProps<any>>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<style scoped></style>
