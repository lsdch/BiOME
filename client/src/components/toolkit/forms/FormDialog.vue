<template>
  <v-dialog
    persistent
    scrollable
    v-model="dialog"
    v-bind="$attrs"
    :max-width="maxWidth ?? 1000"
    :fullscreen="fullscreen || xs"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-card flat :rounded="false">
      <v-toolbar dark dense flat class="position-sticky">
        <v-toolbar-title class="font-weight-bold"> {{ title }} </v-toolbar-title>
        <template #append>
          <v-btn
            color="primary"
            type="submit"
            @click="emit('submit')"
            :loading="loading"
            text="Submit"
          />
          <v-btn color="grey" @click="close" text="Cancel" />
        </template>
      </v-toolbar>
      <v-card-text>
        <!-- Default form slot -->
        <slot></slot>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import { onBeforeMount, useSlots } from 'vue'
import { useDisplay } from 'vuetify'
import { VDialog } from 'vuetify/components'

const dialog = defineModel<boolean>()

const emit = defineEmits<{ submit: []; close: [] }>()

const { xs } = useDisplay()

defineProps<{
  title: string
  loading?: boolean
  fullscreen?: boolean
  maxWidth?: number
}>()

function close() {
  dialog.value = false
  emit('close')
}

const slots = useSlots()
onBeforeMount(() => {
  if (!slots.default) console.error('No content provided in FormDialog slot.')
})
</script>

<style scoped></style>
