<template>
  <v-menu :close-on-content-click="false" origin="left center">
    <template #activator="{ props }">
      <v-input v-bind="$attrs">
        <div class="color-picker-preview">
          <v-avatar :color="model" v-bind="props"> </v-avatar>
        </div>
        <v-label
          v-if="label !== undefined"
          :text="label"
          class="mx-3 opacity-100"
          v-bind="props"
        ></v-label>
      </v-input>
    </template>
    <v-card>
      <v-color-picker v-model="model" v-bind="props" />
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { VColorPicker } from 'vuetify/components'

const model = defineModel<string>()
const { label, ...props } = defineProps<
  {
    label?: string
  } & Omit<VColorPicker['$props'], 'modelValue' | 'onUpdate:modelValue'>
>()
</script>

<style scoped lang="scss">
.color-picker-preview {
  border-radius: 50%;
  width: fit-content;
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAYAAADED76LAAAAAXNSR0IArs4c6QAAACRJREFUKFNjPHTo0H8GJGBnZ8eIzGekgwJk+0BsdCtRHEQbBQBbbh0dIGKknQAAAABJRU5ErkJggg==)
    repeat;
}
</style>
