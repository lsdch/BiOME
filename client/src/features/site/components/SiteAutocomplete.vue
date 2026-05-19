<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchValue"
    :items="items ?? []"
    item-title="name"
    item-value="code"
    :loading="isFetching ? 'primary' : false"
    no-filter
    chips
    clear-on-select
    clearable
    return-object
    label="Site search"
    placeholder="Match by name, code or locality"
    persistent-placeholder
    v-bind="$attrs"
  >
    <template #item="{ item: site, props }">
      <v-list-item :title="site.name" v-bind="props" :disabled="disabledCodes?.includes(site.code)">
        <template #subtitle>
          <v-list-item-subtitle>
            {{ site.locality ?? 'Unspecified locality' }}
            <CountryChip v-if="site.country" :country="site.country" size="small" />
          </v-list-item-subtitle>
        </template>
        <template #append>
          <div class="d-flex flex-column align-end">
            <v-chip :text="site.code" class="font-monospace" size="small" />
          </div>
        </template>
      </v-list-item>
    </template>
    <template #no-data>
      <v-list-item
        class="text-muted"
        :title="
          isFetching
            ? 'Loading...'
            : (searchValue?.length ?? 0) > 2
              ? 'No matching sites'
              : 'Waiting for query...'
        "
      >
        <template #prepend v-if="isFetching">
          <v-progress-circular indeterminate color="primary" class="mr-3" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import { searchSitesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import CountryChip from './CountryChip'

const model = defineModel<SiteItem>()
const searchValue = ref<string | undefined>('')
const { disabledCodes } = defineProps<{
  disabledCodes?: string[]
}>()

const { data: items, isFetching } = useQuery(
  computed(() => ({
    enabled: (searchValue.value?.length ?? 0) > 2,
    ...searchSitesOptions({ query: { query: searchValue.value } })
  }))
)
</script>

<script lang="ts">
/**
 * Provides an autocomplete input for searching sites by name, locality, code or WGS84 coordinates.
 */
export default {}
</script>

<style scoped lang="scss"></style>
