<template>
  <v-dialog
    persistent
    scrollable
    v-model="dialog"
    v-bind="$attrs"
    :max-width="maxWidth ?? 1000"
    :fullscreen="fullscreen ?? xs"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-card flat :rounded="false">
      <v-toolbar class="position-sticky">
        <template #title>
          <slot name="title">
            <v-toolbar-title class="font-weight-bold"> {{ title }} </v-toolbar-title>
          </slot>
        </template>
        <template #append>
          <slot name="append" />
          <v-btn class="ml-1" color="grey" @click="close" :text="closeText" />
        </template>
      </v-toolbar>
      <slot>
        <v-card-text>
          <!-- Default form slot -->
          <slot name="text" />
        </v-card-text>
      </slot>
      <!-- <v-divider></v-divider> -->
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { onBeforeMount, useSlots } from 'vue'
import { useDisplay } from 'vuetify'
import { VDialog } from 'vuetify/components'

const { xs } = useDisplay()
const dialog = defineModel<boolean>({ default: false })

const emit = defineEmits<{ close: [] }>()

export type CardDialogProps = {
  title?: string
  loading?: boolean
  fullscreen?: boolean
  maxWidth?: number
  closeText?: string
}

withDefaults(defineProps<CardDialogProps>(), { closeText: 'Close' })

function close() {
  dialog.value = false
  emit('close')
}

const slots = useSlots()
onBeforeMount(() => {
  if (!slots.default) console.error('No content provided in CardDialog slot.')
})
</script>

<style scoped></style>
