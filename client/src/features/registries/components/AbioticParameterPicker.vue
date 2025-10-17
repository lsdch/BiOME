<template>
  <v-select
    label="Parameter"
    :items
    :loading="isPending"
    v-model="model"
    item-title="label"
    v-bind="$attrs"
    :error-messages="error?.detail"
  >
    <template #item="{ item, props }">
      <v-list-item v-bind="props" :subtitle="item.raw.description">
        <template #prepend="props">
          <slot name="prepend" :props :item></slot>
        </template>
        <template #append>
          <v-chip :text="item.raw.unit" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { AbioticParameter } from '@/api'
import { listAbioticParametersOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'

const model = ref<AbioticParameter>()

const { data: items, error, isPending } = useQuery(listAbioticParametersOptions())
</script>

<style lang="scss" scoped></style>
