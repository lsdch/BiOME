<template>
  <v-card id="panels">
    <v-tabs v-model="tab" :mandatory="false">
      <v-tab :disabled value="occurrences" prepend-icon="mdi-crosshairs-gps">
        Occurrences
        <template #append>
          <v-badge
            color="primary"
            inline
            :content="
              site?.samplings?.reduce(
                (acc, { occurrences }) => acc + (occurrences?.length ?? 0),
                0
              ) ?? 0
            "
          />
        </template>
      </v-tab>
      <v-tab :disabled value="samplings" prepend-icon="mdi-package-down">
        Samplings
        <template #append>
          <v-badge color="primary" inline :content="site?.samplings?.length ?? 0" />
        </template>
      </v-tab>
      <v-tab :disabled value="abiotics" prepend-icon="mdi-gauge">
        Abiotic
        <template #append>
          <v-badge
            :color="site?.abiotic_measurements?.length ? 'primary' : ''"
            inline
            :content="site?.abiotic_measurements?.length ?? 0"
          />
        </template>
      </v-tab>
      <v-tab :disabled value="datasets">
        Datasets
        <template #append>
          <v-badge color="purple" inline :content="site?.datasets?.length ?? 0" variant="tonal" />
        </template>
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="tab" crossfade direction="vertical">
      <v-tabs-window-item value="occurrences">
        <OccurrencesAtSiteTable :samplings="site?.samplings ?? []" v-if="site" />
      </v-tabs-window-item>

      <v-tabs-window-item value="samplings">
        <SamplingTableAtSite :samplings="site?.samplings ?? []"> </SamplingTableAtSite>
      </v-tabs-window-item>

      <v-tabs-window-item value="abiotics">
        <v-card-text>
          <AbioticTab :abiotic_measurements="site?.abiotic_measurements ?? []" />
        </v-card-text>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup lang="ts">
import { Site } from '@/api'
import OccurrencesAtSiteTable from '@/features/occurrences/components/OccurrencesAtSiteTable.vue'
import SamplingTableAtSite from '@/features/site/components/SamplingTableAtSite.vue'
import AbioticTab from '@/features/site/components/AbioticTab.vue'
import { computed, ref } from 'vue'

const { site, ...props } = defineProps<{
  site?: Site
  disabled?: boolean
}>()

const disabled = computed(() => props.disabled || !site)

type Tabs = 'occurrences' | 'samplings' | 'abiotics' | 'datasets'
const tab = ref<Tabs>()
</script>

<style scoped lang="scss"></style>
