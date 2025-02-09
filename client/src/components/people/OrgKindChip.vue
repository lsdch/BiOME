<template>
  <v-chip
    :title="label"
    label
    variant="outlined"
    :class="[$attrs['class'], 'px-1']"
    v-bind="{ ...OrgKind.props[kind], ...$attrs }"
  >
    <template v-slot:prepend>
      <v-icon v-bind="OrgKind.props[kind]" size="small"></v-icon>
    </template>
    <template v-slot v-if="!hideLabel">
      <slot>
        {{ label }}
      </slot>
    </template>
  </v-chip>
</template>

<script setup lang="ts">
import { OrgKind } from '@/api/adapters'

defineProps<{
  kind: OrgKind
  label?: string
  hideLabel?: boolean
}>()
</script>

<style scoped></style>
