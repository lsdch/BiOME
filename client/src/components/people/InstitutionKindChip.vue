<template>
  <v-chip
    :title="label"
    label
    variant="outlined"
    :class="[$attrs['class'], 'px-1']"
    v-bind="{ ...InstitutionKind.props[kind], ...$attrs }"
  >
    <template v-slot:prepend>
      <v-icon v-bind="InstitutionKind.props[kind]" size="small"></v-icon>
    </template>
    <template v-slot v-if="!hideLabel">
      <slot>
        {{ label }}
      </slot>
    </template>
  </v-chip>
</template>

<script setup lang="ts">
import { InstitutionKind } from '@/api/adapters'

defineProps<{
  kind: InstitutionKind
  label?: string
  hideLabel?: boolean
}>()
</script>

<style scoped></style>
