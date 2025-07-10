<template>
  <v-autocomplete v-model="model" :items="presets" item-title="name" return-object v-bind="$attrs">
    <template #item="{ item, props }">
      <v-list-item :title="item.raw.name" v-bind="props">
        <template #subtitle> by {{ item.raw.meta.created_by?.name }} </template>
        <template #append>
          <div class="d-flex ga-3 align-center">
            <MapPresetSummaryIcons :spec="item.raw.spec" />
          </div>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { listMapPresetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { pluralize } from '@/functions/text'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { ParsedMapPreset, parseMapPreset } from './map-presets'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import MapPresetSummaryIcons from './MapPresetSummaryIcons.vue'

const model = defineModel<ParsedMapPreset>()

const { include } = defineProps<{
  include: ('personal' | 'globals' | 'public')[]
}>()

const { data } = useQuery(listMapPresetsOptions())

const { user } = storeToRefs(useUserStore())

const presets = computed<ParsedMapPreset[] | undefined>(() => {
  return data.value?.reduce<ParsedMapPreset[]>((acc, v) => {
    const p = parseMapPreset(v)
    if (
      (include.includes('personal') && p.meta.created_by?.id == user.value?.id) ||
      (include.includes('globals') && p.is_global) ||
      (include.includes('public') && p.is_public)
    ) {
      acc.push(p)
    }
    return acc
  }, [])
})
</script>

<style scoped lang="scss"></style>
