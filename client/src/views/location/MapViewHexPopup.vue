<template>
  <v-card :min-width="200">
    <v-list density="compact">
      <SiteListDialog :sites="data?.flatMap(({ data }) => data)" :max-width="1200">
        <template #activator="{ props }">
          <v-list-item v-bind="props">
            {{ pluralizeWithCount(data?.length ?? 0, 'site') }}
          </v-list-item>
        </template>
      </SiteListDialog>
      <OccurrenceListDialog
        with-site
        :occurrences="
          data?.flatMap(({ data: { occurrences, ...site } }) => {
            return occurrences.map((o) => ({
              ...o,
              site
            }))
          })
        "
        :max-width="1200"
      >
        <template #activator="{ props }">
          <v-list-item v-bind="props">
            {{
              pluralizeWithCount(
                data?.map(({ data }) => data.occurrences.length).reduce((a, b) => a + b, 0) ?? 0,
                'occurrence'
              )
            }}
          </v-list-item>
        </template>
      </OccurrenceListDialog>
      <AreaSampledTaxa :data />
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import { HexPopupData } from '@/components/maps/SitesMap.vue'
import AreaSampledTaxa from './AreaSampledTaxa.vue'
import { pluralizeWithCount } from '@/functions/text'
import OccurrenceListDialog from './OccurrenceListDialog.vue'
import SiteListDialog from './SiteListDialog.vue'

defineProps<{ data: HexPopupData<SiteWithOccurrences>[] | undefined }>()
</script>

<style scoped lang="scss"></style>
