<template>
  <v-dialog
    v-model="dialog"
    :max-width="maxWidth ?? 1000"
    :fullscreen="fullscreen || $vuetify.display.xs"
    persistent
    scrollable
    :activator
    v-bind="$attrs"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData" v-if="slots.activator">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-card flat :rounded="false" :title :subtitle class="overflow-x-auto" :prepend-icon>
      <template
        v-for="name in slotNames.filter((s) => !['activator', 'append'].includes(s))"
        #[name]="slotData"
      >
        <slot :name v-bind="(slotData as any) ?? {}" />
      </template>
      <template #append>
        <slot name="append" />
        <v-btn
          class="ml-1"
          color="grey"
          @click="close"
          v-bind="
            $vuetify.display.smAndUp
              ? {
                  text: closeText
                }
              : {
                  icon: 'mdi-close',
                  size: 'small'
                }
          "
          variant="plain"
        />
      </template>
      <v-divider />
      <slot>
        <v-card-text>
          <!-- Default form slot -->
          <slot name="text" />
        </v-card-text>
      </slot>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ComponentPublicInstance } from 'vue'
import { VCard, VDialog } from 'vuetify/components'

const dialog = defineModel<boolean>()

const emit = defineEmits<{ close: [] }>()

export type CardDialogProps = {
  title?: string
  subtitle?: string
  loading?: boolean
  fullscreen?: boolean
  maxWidth?: number
  closeText?: string
  prependIcon?: IconValue
  activator?: (string & {}) | Element | 'parent' | ComponentPublicInstance
}

withDefaults(defineProps<CardDialogProps>(), { closeText: 'Close' })

function close() {
  dialog.value = false
  emit('close')
}

// type SlotType = VCard['$slots'] & Pick<VDialog['$slots'], 'activator'>
const slots = defineSlots<VCard['$slots'] & Pick<VDialog['$slots'], 'activator'>>()
const slotNames = Object.keys(slots) as Array<keyof typeof slots>
</script>

<style scoped></style>
