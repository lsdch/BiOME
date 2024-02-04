<template>
  <v-dialog v-model="dialog" max-width="500px">
    <v-card>
      <v-toolbar dark dense flat>
        <v-toolbar-title class="text-body-2 font-weight-bold grey--text">
          {{ content.title }}
        </v-toolbar-title>
      </v-toolbar>
      <v-card-text> {{ content.message }} </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="grey" variant="text" @click="cancel" text="Cancel" />
        <v-btn color="blue-darken-1" variant="text" @click="agree(content.payload)" text="OK" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="Payload = any">
export type ConfirmDialogProps<Payload> = {
  title: string
  message: string
  payload?: Payload
}

const dialog = defineModel<boolean>()

const content = defineProps<ConfirmDialogProps<Payload>>()

const emit = defineEmits<{
  agree: [payload?: Payload]
  cancel: []
}>()

function agree(payload?: Payload) {
  emit('agree', payload)
  dialog.value = false
}

function cancel() {
  emit('cancel')
  dialog.value = false
}
</script>

<style scoped></style>
