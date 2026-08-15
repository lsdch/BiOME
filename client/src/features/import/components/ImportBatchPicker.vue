<template>
  <v-autocomplete
    v-model="model"
    :label="multiple ? 'Import batches' : 'Import batch'"
    :multiple
    :return-object
    :items="batches"
    :loading
    item-title="label"
    no-data-text="No import batches found"
  ></v-autocomplete>
</template>

<script setup lang="ts" generic="Multiple extends boolean, ReturnObject extends boolean">
import { ImportBatch } from '@/api'
import { listImportBatchesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'

const model = defineModel<Value<ImportBatch, ReturnObject, Multiple>>()

const { multiple } = defineProps<{
  multiple?: Multiple
  returnObject?: ReturnObject
}>()

const { data: batches, isPending: loading } = useQuery(listImportBatchesOptions())
</script>

<style scoped lang="scss"></style>
