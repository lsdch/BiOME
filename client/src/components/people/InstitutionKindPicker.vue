<template>
  <v-select v-model="model" :items v-bind="$attrs">
    <template v-if="model" #prepend-inner>
      <v-icon v-bind="InstitutionKind.props[model]" />
    </template>
    <template #item="{ item, props }">
      <v-list-item :title="item.title" v-bind="props">
        <template #prepend>
          <v-icon v-bind="InstitutionKind.props[item.raw.value]" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { $InstitutionKind } from '@/api'
import { InstitutionKind } from '@/api/adapters'

const model = defineModel<InstitutionKind>()
const items = $InstitutionKind.enum.map((value) => ({
  value,
  title: InstitutionKind.humanize(value)
}))
</script>

<style scoped lang="scss"></style>
