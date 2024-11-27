<template>
  <v-form>
    <TaxonPicker
      label="Taxa"
      v-model="model.target_taxa"
      item-value="name"
      multiple
      chips
      closable-chips
      clearable
    />
    <v-text-field label="Comments" v-model="model.comments" />
  </v-form>
</template>

<script setup lang="ts">
import { Spotting, SpottingUpdate } from '@/api'
import { FormProps, useForm } from '../toolkit/forms/form'
import { ref } from 'vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'

const sampling = defineModel<Spotting>('sampling', { required: true })

const model = ref<SpottingUpdate>({
  target_taxa: sampling.value.target_taxa.map(({ name }) => name),
  comments: sampling.value.comments
})
</script>

<style scoped lang="scss"></style>
