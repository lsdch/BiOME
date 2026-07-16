<template>
  <TaxonPicker v-model="model" v-bind="props">
    <template #append-inner="props">
      <slot name="append-inner" v-bind="props" />
      <v-btn :active="wholeclade" icon="mdi-family-tree" size="x-small" @click.stop="wholeclade = !wholeClade"
        :variant="wholeClade ? 'outlined' : 'plain'" :color="wholeClade ? 'primary' : ''" v-tooltip="{
          location: 'top',
          text: 'When active, include descendant taxa'
        }"></v-btn>
    </template>
  </TaxonPicker>
</template>

<script setup lang="ts">
import { TaxaFilters } from '@/api'
import TaxonPicker, { TaxonPickerProps } from './TaxonPicker.vue'

const model = defineModel<TaxaFilters['taxa']>('taxa')
const wholeclade = defineModel<TaxaFilters['whole_clade']>('wholeClade')
const props = defineProps<TaxonPickerProps<'name' | 'id', true, false>>()
</script>

<style scoped lang="scss"></style>
