<template>
  <v-dialog v-model="dialog" max-width="500px">
    <v-card>
      <v-toolbar dark dense flat>
        <v-toolbar-title class="text-body-2 font-weight-bold grey--text">
          {{ props.title }}
        </v-toolbar-title>
      </v-toolbar>
      <v-card-text> {{ props.message }} </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="grey" variant="text" @click="cancel" text="Cancel" />
        <v-btn color="blue-darken-1" variant="text" @click="confirm" text="OK" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="Payload = any">
export type ConfirmDialogProps = {
  title: string
  message: string
  onConfirm?: () => any
  onCancel?: () => any
}

const dialog = defineModel<boolean>()

const props = defineProps<ConfirmDialogProps>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

function confirm() {
  emit('confirm')
  dialog.value = false
  return props.onConfirm?.()
}

function cancel() {
  emit('cancel')
  dialog.value = false
  return props.onCancel?.()
}
</script>

<style scoped></style>
