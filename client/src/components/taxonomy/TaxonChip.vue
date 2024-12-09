<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip :text="taxon.name" v-bind="{ ...props, ...$attrs }"> </v-chip>
    </template>
    <v-card
      :title="taxon.name"
      :subtitle="taxon.authorship"
      class="bg-surface-light"
      density="compact"
      :to="{ name: 'taxonomy', hash: `#${taxon.name}` }"
    >
      <template #prepend>
        <LinkIconGBIF
          v-if="taxon.GBIF_ID"
          :GBIF_ID="taxon.GBIF_ID"
          variant="tonal"
          size="x-small"
          @click.stop
        />
        <FTaxonStatusIndicator v-else :status="taxon.status" />
      </template>
      <!-- <template #append>
            <v-btn
              icon="mdi-link-variant"
              :to="{ name: 'taxonomy', hash: `#${taxon.name}` }"
              variant="plain"
            />
          </template>
        </v-card-item>
      </template> -->

      <v-card-text>
        <div class="d-flex justify-space-between">
          <v-chip :text="taxon.status" class="ma-1" />
          <v-chip :text="taxon.rank" class="ma-1" />
        </div>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Taxon } from '@/api'
import LinkIconGBIF from './LinkIconGBIF.vue'
import { FTaxonStatusIndicator } from './functionals'

const props = defineProps<{ taxon: Taxon }>()
</script>

<style scoped lang="scss"></style>
