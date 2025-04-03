<template>
  <v-card class="fill-height">
    <v-toolbar density="compact" color="white">
      <v-toolbar-title :text="name" />
      <v-spacer />
      <LinkIconGBIF :GBIF_ID="GBIF_ID" />
    </v-toolbar>
    <v-card-subtitle>{{ rank }}</v-card-subtitle>

    <div class="d-flex justify-space-between mb-5 align-center">
      <div>
        <v-card-subtitle> Imported {{ imported }} nodes </v-card-subtitle>
        <v-card-subtitle class="flex-end"> Started {{ elapsed }} </v-card-subtitle>
      </div>
      <div class="mr-5">
        <v-tooltip v-if="error" location="left">
          <template v-slot:activator="{ props }">
            <v-icon color="error" v-bind="props" icon="mdi-alert" />
          </template>
          An unexpected error occurred during the import.
        </v-tooltip>
        <v-icon v-else-if="done" color="green" icon="mdi-check" />
        <v-progress-circular v-else indeterminate color="blue" />
      </div>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import LinkIconGBIF from '@/components/taxonomy/LinkIconGBIF'

export type ImportProcess = {
  name: string
  GBIF_ID: number
  expected: number
  imported: number
  rank: string
  started: string
  done: boolean
  elapsed?: string
  error?: object
}

defineProps<ImportProcess>()
</script>

<style scoped></style>
