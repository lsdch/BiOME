<template>
  <TaxonPicker
    v-model="model"
    @update:model-value="(v) => console.log('update', v)"
    v-bind="props"
    item-value="id"
    :return-object="false"
  >
    <template #append-inner="props">
      <slot name="append-inner" v-bind="props" />
      <v-btn
        :active="wholeclade"
        icon="mdi-family-tree"
        size="x-small"
        @click.stop="wholeclade = !wholeClade"
        :variant="wholeClade ? 'outlined' : 'plain'"
        :color="wholeClade ? 'primary' : ''"
        v-tooltip="{
          location: 'top',
          text: 'When active, include descendant taxa'
        }"
      ></v-btn>
    </template>
  </TaxonPicker>
</template>

<script setup lang="ts">
import TaxonPicker, { TaxonPickerProps } from './TaxonPicker.vue'

const model = defineModel<UUID[]>('taxa')
const wholeclade = defineModel<boolean>('wholeClade')
const props = defineProps<TaxonPickerProps<true, false>>()
</script>

<style scoped lang="scss"></style>
