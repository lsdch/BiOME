<template>
  <v-switch
    :label
    v-model="model"
    :indeterminate="model === undefined"
    :color="colorTrue"
    :base-color="model === undefined ? baseColor : colorFalse"
    v-bind="$attrs"
  >
    <template #append>
      <v-icon v-if="model !== undefined" icon="mdi-close-circle" @click="model = undefined" />
    </template>
    <template #message="{ message }">
      <slot v-if="$slots['message']" name="message" :message />
      <template v-else>{{ message }}</template>
    </template>
  </v-switch>
</template>

<script setup lang="ts">
import { VSwitch } from 'vuetify/components'

const model = defineModel<boolean | undefined>()

const props = defineProps<{
  label?: string
  colorTrue?: string
  baseColor?: string
  colorFalse?: string
}>()
</script>

<style scoped lang="scss"></style>
