<template>
  <v-chip
    :text="`${db.code}:${accession}`"
    :href="SeqReference.link(seqRef)"
    target="_blank"
    v-bind="$attrs"
    :color="is_origin ? 'success' : undefined"
  />
</template>

<script setup lang="ts">
import { SeqReference } from '@/api/adapters'

const { seqRef } = defineProps<{ seqRef: SeqReference }>()
const { accession, db, is_origin } = seqRef
</script>

<style scoped lang="scss"></style>
