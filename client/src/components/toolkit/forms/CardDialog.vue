<template>
  <v-dialog
    persistent
    scrollable
    v-model="dialog"
    :max-width="maxWidth ?? 1000"
    :fullscreen="fullscreen ?? xs"
    v-bind="$attrs"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-card flat :rounded="false" :title :subtitle class="overflow-x-auto">
      <!-- <v-toolbar class="position-sticky">
        <template #title>
          <slot name="title">
            <v-toolbar-title class="font-weight-bold"> {{ title }} </v-toolbar-title>
          </slot>
        </template>
      </v-toolbar> -->
      <template
        v-for="(name, index) of Object.keys($slots).filter((k) => k !== 'append')"
        #[name]="slotData"
        :key="index"
      >
        <slot :name v-bind="slotData ?? {}" />
      </template>
      <template #append>
        <slot name="append" />
        <v-btn class="ml-1" color="grey" @click="close" :text="closeText" variant="plain" />
      </template>
      <v-divider />
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
import { useDisplay } from 'vuetify'
import { VDialog } from 'vuetify/components'

const { xs } = useDisplay()
const dialog = defineModel<boolean>({ default: false })

const emit = defineEmits<{ close: [] }>()

export type CardDialogProps = {
  title?: string
  subtitle?: string
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
</script>

<style scoped></style>
