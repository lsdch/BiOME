<template>
  <v-list>
    <v-list-item>
      <DatasetPicker
        v-model="context.datasets"
        label="Datasets"
        placeholder="None"
        persistent-placeholder
        density="compact"
        clearable
        multiple
        chips
        closable-chips
        hide-details
        persistent-hint
        class="my-3"
      />
    </v-list-item>
    <v-list-item>
      <TaxonPicker
        class="mt-1"
        v-model="context.taxa"
        item-value="name"
        density="compact"
        multiple
        chips
        closable-chips
        hide-details
      >
      </TaxonPicker>
    </v-list-item>
    <v-list-item>
      <v-switch
        class="px-2"
        label="Use clade"
        v-model="context.whole_clade"
        color="primary"
        density="compact"
        :disabled="!context.taxa?.length"
        hide-details
      >
        <template #append>
          <InlineHelp>
            When enabled, all occurrences of descendant taxa will be included.
          </InlineHelp>
        </template>
      </v-switch>
    </v-list-item>
    <div class="d-flex justify-end ga-1 mx-2">
      <v-btn
        prepend-icon="mdi-close"
        text="Clear"
        rounded="sm"
        variant="plain"
        color=""
        @click="context = {}"
      />
      <v-btn
        rounded="sm"
        variant="tonal"
        prepend-icon="mdi-tray-arrow-down"
        text="Apply"
        v-tooltip="'Apply context filters to all data feeds'"
        @click="applyContext()"
      />
    </div>
  </v-list>
</template>

<script setup lang="ts">
import DatasetPicker from '../datasets/DatasetPicker.vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import InlineHelp from '../toolkit/ui/InlineHelp.vue'
import { useDataFeeds } from './data_feeds'

const { context, applyContext } = useDataFeeds()
</script>

<style scoped lang="scss"></style>
