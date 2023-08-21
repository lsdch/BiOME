<template>
  <v-card class="h-100 d-flex flex-column">
    <v-toolbar density="compact" color="white">
      <v-toolbar-title :title="name">
        {{ name }}
      </v-toolbar-title>
      <v-tooltip location="top" text="Drop anchor">
        <template v-slot:activator="{ isActive, props }">
          <v-btn
            v-bind="props"
            :icon="isActive ? 'mdi-close' : 'mdi-anchor'"
            size="small"
            :color="isActive ? 'red' : 'orange'"
          ></v-btn>
        </template>
      </v-tooltip>
      <LinkIconGBIF :GBIF_ID="GBIF_ID" />
    </v-toolbar>
    <v-card-subtitle class="d-flex flex-column flex-md-row">
      <span class="col-12">{{ rank }}</span> <v-spacer></v-spacer> <span>{{ authorship }}</span>
    </v-card-subtitle>
    <div class="mb-5 d-flex flex-column flex-grow-1 justify-end">
      <v-card-subtitle class="flex-end">
        Last modified
        {{ moment(modified ?? created).format('DD MMM y HH:MM') }}
      </v-card-subtitle>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import LinkIconGBIF from './LinkIconGBIF.vue'
import moment from 'moment'

import type { Taxon } from '@/types/taxonomy'

defineProps<Taxon>()
</script>

<style scoped></style>
