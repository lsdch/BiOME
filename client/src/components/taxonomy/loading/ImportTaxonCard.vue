<template>
  <v-card>
    <v-toolbar density="compact" color="white">
      <v-toolbar-title>
        {{ name }}
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <LinkIconGBIF :GBIF_ID="GBIF_ID" />
    </v-toolbar>
    <v-card-subtitle>{{ rank }}</v-card-subtitle>

    <div class="d-flex justify-space-between mb-5 align-center">
      <div>
        <v-card-subtitle> Imported {{ imported }} nodes </v-card-subtitle>
        <v-card-subtitle class="flex-end"> Started {{ elapsed }} </v-card-subtitle>
      </div>
      <v-icon v-if="done" class="mr-5" color="green">mdi-check</v-icon>
      <v-progress-circular v-else indeterminate class="mr-5" color="blue"/>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import LinkIconGBIF from '@/components/taxonomy/LinkIconGBIF.vue'

export type ImportProcess = {
  name: string
  GBIF_ID: number
  expected: number
  imported: number
  rank: string
  started: string
  done: boolean
  elapsed?: string
}

defineProps<ImportProcess>()
</script>

<style scoped></style>
