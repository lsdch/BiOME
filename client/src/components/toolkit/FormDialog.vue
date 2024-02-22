<template>
  <v-dialog persistent v-model="dialog" v-bind="$attrs" :fullscreen="xs" max-width="1000">
    <!-- Expose activator slot -->
    <template v-slot:activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-card flat :rounded="false">
      <v-toolbar dark dense flat>
        <v-toolbar-title class="font-weight-bold"> {{ title }} </v-toolbar-title>
        <template v-slot:append>
          <v-btn color="grey" @click="dialog = false">Cancel</v-btn>
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
import { computed, onBeforeMount, useSlots } from 'vue'
import { useDisplay } from 'vuetify'
import { VDialog } from 'vuetify/components'
import { Mode } from './form'

const dialog = defineModel<boolean>()

const { xs } = useDisplay()

const props = defineProps<{
  entityName: string
  mode: Mode
}>()

const title = computed(() => {
  return `${props.mode} ${props.entityName}`
})

const slots = useSlots()
onBeforeMount(() => {
  if (!slots.default) console.error('No content provided in FormDialog slot.')
})
</script>

<style scoped></style>
