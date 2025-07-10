<template>
  <CRUDTable
    :fetch-items
    :fetch-params
    :headers
    entity-name="Map presets"
    :delete="{
      mutation: (opts?: Partial<Options<DeleteMapPresetData>>) =>
        deleteMapPresetMutation({
          ...opts,
          responseTransformer: async (data) => {
            return await deleteMapPresetResponseTransformer(parseMapPreset(data as MapToolPreset))
          }
        }) as UseMutationOptions<
          ParsedMapPreset,
          DeleteMapPresetError,
          Options<DeleteMapPresetData>
        >,
      params: ({ name }) => ({ path: { name } })
    }"
  >
    <template #item.name="{ item }: { item: ParsedMapPreset }">
      <ActivableField
        v-slot="{ actions, proxy, isPristine, active, props, cancel, save }"
        :activable="true"
        v-model="item.name"
      >
        <v-text-field
          v-model="proxy.value"
          hide-details
          density="compact"
          :variant="active ? undefined : 'plain'"
          :readonly="!active"
        >
          <template #append-inner>
            <template v-if="active">
              <v-btn
                icon="mdi-check"
                size="small"
                variant="plain"
                @click="updateItem(item, { name: proxy.value }, cancel).then(save)"
              />
              <v-btn icon="mdi-close" size="small" variant="plain" @click="cancel()" />
            </template>
            <v-btn
              v-else
              icon="mdi-pencil"
              size="small"
              density="compact"
              variant="plain"
              @click="props.onfocus()"
            />
          </template>
        </v-text-field>
      </ActivableField>
    </template>
    <template #item.overview="{ item }: { item: ParsedMapPreset }">
      <div class="d-flex ga-5 justify-center">
        <MapPresetSummaryIcons :spec="item.spec" />
      </div>
    </template>
    <template #item.is_public="{ item }: { item: ParsedMapPreset }">
      <v-confirm-edit v-model="item.is_public" v-slot="{ actions, model: proxy, cancel, save }">
        <v-switch
          v-model="proxy.value"
          @update:model-value="
            (v) => updateItem(item, { is_public: proxy.value }, cancel).then(save)
          "
          color="success"
          hide-details
        />
      </v-confirm-edit>
    </template>
    <template #expanded-row-inject="{ item }: { item: ParsedMapPreset }">
      <ActivableField
        v-slot="{ actions, proxy, isPristine, active, props, cancel, save }"
        activable
        v-model="item.description"
      >
        <v-textarea
          class="mx-5 my-3"
          v-model="proxy.value"
          label="Description"
          hide-details
          density="compact"
          :variant="active ? undefined : 'plain'"
          :readonly="!active"
          auto-grow
          :rows="2"
        >
          <template #prepend>
            <div v-if="active" class="d-flex flex-column">
              <v-btn
                icon="mdi-check"
                size="small"
                variant="plain"
                @click="updateItem(item, { description: proxy.value }, cancel).then(save)"
              />
              <v-btn icon="mdi-close" size="small" variant="plain" @click="cancel()" />
            </div>
            <v-btn
              v-else
              icon="mdi-pencil"
              size="small"
              density="compact"
              variant="plain"
              @click="props.onfocus()"
            />
          </template>
        </v-textarea>
      </ActivableField>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { DeleteMapPresetData, DeleteMapPresetError, ErrorModel, MapToolPreset } from '@/api'
import {
  createUpdateMapPresetMutation,
  deleteMapPresetMutation,
  listMapPresetsOptions,
  listMapPresetsQueryKey
} from '@/api/gen/@tanstack/vue-query.gen'
import { Options } from '@hey-api/client-fetch'
import { UndefinedInitialQueryOptions, useMutation, UseMutationOptions } from '@tanstack/vue-query'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import { ParsedMapPreset, parseMapPreset } from './map-presets'
import {
  deleteMapPresetResponseTransformer,
  listMapPresetsResponseTransformer
} from '@/api/gen/transformers.gen'
import ActivableField from '../toolkit/forms/ActivableField.vue'
import MapPresetSummaryIcons from './MapPresetSummaryIcons.vue'
import { useFeedback } from '@/stores/feedback'
import { computed } from 'vue'
import { computedWithControl } from '@vueuse/core'

const props = defineProps<{ all?: boolean }>()

const showAllRef = computed(() => props.all ?? false)

const fetchItems = computedWithControl(showAllRef, () => {
  console.log('triggered')
  return {
    ...(listMapPresetsOptions({
      query: { all: showAllRef.value },
      responseTransformer: async (data) =>
        await listMapPresetsResponseTransformer((data as MapToolPreset[]).map(parseMapPreset))
    }) as UndefinedInitialQueryOptions<ParsedMapPreset[], ErrorModel, ParsedMapPreset[]>),
    queryKey: [...listMapPresetsQueryKey({ query: { all: showAllRef.value } }), 'transformed']
  }
})

const fetchParams = computed(() => ({
  query: { all: showAllRef.value }
}))

const headers: CRUDTableHeader<ParsedMapPreset>[] = [
  {
    key: 'name',
    title: 'Name',
    sortable: true,
    align: 'start'
  },
  { key: 'overview', title: '' },
  { key: 'is_public', title: 'Public', width: 0, align: 'center' }
]

const { mutateAsync } = useMutation(createUpdateMapPresetMutation())

const { feedback } = useFeedback()

function updateItem(
  { meta, $schema, ...item }: ParsedMapPreset,
  update: Partial<ParsedMapPreset>,
  onError?: (error: ErrorModel) => void
) {
  return mutateAsync({
    body: { ...item, ...update, spec: JSON.stringify(item.spec) }
  }).catch((error: ErrorModel) => {
    onError?.(error)
    feedback({
      type: 'error',
      message: `Failed to update map preset "${item.name}": ${error.detail}`
    })
  })
}
</script>

<style scoped lang="scss"></style>
