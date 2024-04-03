<template>
  <v-chip label variant="outlined" v-bind="{ ...kindIcon(props.kind), ...$attrs }">
    <template v-slot:prepend>
      <v-icon v-bind="kindIcon(props.kind)"></v-icon>
    </template>
    <template v-slot>
      <slot></slot>
    </template>
  </v-chip>
</template>

<script setup lang="ts">
import { InstitutionKind } from '@/api'
import { kindIcon } from './institutionKind'

const props = defineProps<{ kind: InstitutionKind }>()
</script>

<style scoped></style>
