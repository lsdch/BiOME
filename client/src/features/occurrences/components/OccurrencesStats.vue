<template>
  <v-list>
    <v-list-item title="Sampling events">
      <template #append>
        <v-badge
          inline
          :content="cells?.reduce((sum, s) => sum + s.samplings_count, 0) ?? 0"
          color="info"
        ></v-badge>
      </template>
    </v-list-item>
    <v-list-item title="Occurrences">
      <template #append>
        <v-badge
          inline
          :content="cells?.reduce((sum, s) => sum + s.occurrences_count, 0) ?? 0"
          color="success"
        />
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { H3CellWithRichness } from '@/api'

const { cells } = defineProps<{
  cells?: H3CellWithRichness[]
}>()
</script>

<style scoped lang="scss"></style>
