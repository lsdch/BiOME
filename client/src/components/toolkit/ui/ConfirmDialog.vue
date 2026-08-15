<template>
  <v-dialog v-model="isRevealed" :width :max-width persistent @keyup.esc="cancel()">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <v-form>
      <template #default="{ isValid, items }">
        <v-card>
          <v-toolbar dark dense flat>
            <v-toolbar-title class="text-body-2 font-weight-bold text-medium-emphasis">
              {{ title }}
            </v-toolbar-title>
          </v-toolbar>
          <slot>
            <v-card-text v-if="message"> {{ message }} </v-card-text>
          </slot>
          <v-card-actions>
            <v-spacer />
            <v-btn color="grey" variant="text" @click="cancel()" text="Cancel" />
            <v-btn
              color="blue-darken-1"
              variant="text"
              @click="confirm()"
              text="OK"
              :disabled="!!items.length && !isValid"
            />
          </v-card-actions>
        </v-card>
      </template>
    </v-form>
  </v-dialog>
</template>

<script setup lang="ts">
export type ConfirmDialogProps<T> = {
  title: string
  message?: string
  data?: T
  width?: number
  maxWidth?: number
}

const isRevealed = defineModel<boolean>()

const { width = 800 } = defineProps<ConfirmDialogProps<any>>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

function confirm() {
  isRevealed.value = false
  emit('confirm')
}

function cancel() {
  isRevealed.value = false
  emit('cancel')
}
</script>

<style scoped></style>
