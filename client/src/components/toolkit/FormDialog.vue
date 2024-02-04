<template>
  <v-dialog persistent v-model="dialog" v-bind="$attrs" :fullscreen="xs">
    <template v-for="(name, index) of slotNames" v-slot:[name]="slotData" :key="index">
      <slot :name="name" v-bind="slotData || {}" />
    </template>

    <v-card flat :rounded="false">
      <v-toolbar dark dense flat>
        <v-toolbar-title class="font-weight-bold"> {{ title }} </v-toolbar-title>
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
import { computed, type Component } from 'vue'
import { ComponentSlots } from 'vue-component-type-helpers'
import { useDisplay } from 'vuetify'
import { VDialog } from 'vuetify/components'
import type { Emits, Mode, Props } from './form'

const dialog = defineModel<boolean>()

const { xs } = useDisplay()

const slots = defineSlots<ComponentSlots<typeof VDialog>>()
const slotNames = Object.keys(slots) as 'default'[]

const props = defineProps<
  {
    // FIXME : stronger type checking for emit events
    form: Component
    entityName: string
  } & Omit<Props<ItemType>, 'onSuccess'>
>()

const title = computed(() => {
  const mode: Mode = props.edit ? 'Edit' : 'Create'
  return `${mode} ${props.entityName}`
})

const emit = defineEmits<Emits<ItemType>>()

function onSuccess(mode: Mode, item: ItemType) {
  dialog.value = false
  emit('success', mode, item)
}
</script>

<style scoped></style>
