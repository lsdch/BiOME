<template>
  <v-select v-model="model" :items v-bind="$attrs">
    <template v-if="model" #prepend-inner>
      <v-icon v-bind="OrgKind.props[model]" />
    </template>
    <template #item="{ item, props }">
      <v-list-item :title="item.title" v-bind="props">
        <template #prepend>
          <v-icon v-bind="OrgKind.props[item.raw.value]" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { $OrgKind } from '@/api'
import { OrgKind } from '@/api/adapters'

const model = defineModel<OrgKind>()
const items = $OrgKind.enum.map((value) => ({
  value,
  title: OrgKind.humanize(value)
}))
</script>

<style scoped lang="scss"></style>
