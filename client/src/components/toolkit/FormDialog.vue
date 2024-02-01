<template>
  <v-dialog persistent v-model="dialog" v-bind="$attrs">
    <template v-for="(name, index) of slotNames" v-slot:[name]="slotData" :key="index">
      <slot :name="name" v-bind="slotData || {}" />
    </template>

    <v-card>
      <v-toolbar dark dense flat>
        <v-toolbar-title class="font-weight-bold"> Create new </v-toolbar-title>
        <template v-slot:append>
          <v-btn color="grey" @click="dialog = false">Cancel</v-btn>
        </template>
      </v-toolbar>
      <v-card-text>
        <component :is="form" v-bind="$props" :onSuccess="onSuccess" />
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import { ref, watch } from 'vue'
import { ComponentSlots } from 'vue-component-type-helpers'
import { VDialog } from 'vuetify/components'
import { type Component } from 'vue'
import type { Props, Emits, Mode } from './form'
const dialog = ref(false)

const slots = defineSlots<ComponentSlots<typeof VDialog>>()
const slotNames = Object.keys(slots) as 'default'[]

defineProps<
  {
    // FIXME : stronger type checking for emit events
    form: Component
  } & Omit<Props<ItemType>, 'onSuccess'>
>()

const emit = defineEmits<Emits<ItemType>>()

function onSuccess(mode: Mode, item: ItemType) {
  close()
  emit('success', mode, item)
}

function open() {
  dialog.value = true
}

function close() {
  dialog.value = false
}

defineExpose({ open, close })
</script>

<style scoped></style>
