<template>
  <v-card :class>
    <v-row>
      <v-col cols="12" md="6">
        <v-expand-transition>
          <v-card
            prepend-icon="mdi-map-marker"
            title="Site"
            subtitle="Pick or create a site"
            flat
            :rounded="0"
            v-show="!site || showEdit"
          >
            <v-divider />
            <v-card-text>
              <SiteAutocomplete @update:model-value="updateSite" />
            </v-card-text>
            <template #append>
              <v-btn text="New site" prepend-icon="mdi-plus" variant="tonal" rounded="md" />
            </template>
            <template #actions v-if="site">
              <v-spacer />
              <v-btn text="Cancel" @click="toggleEdit(false)"></v-btn>
            </template>
          </v-card>
        </v-expand-transition>
        <v-expand-transition>
          <SitePreviewCard v-show="!showEdit" :site rounded="0" flat @edit="toggleEdit(true)" />
        </v-expand-transition>
      </v-col>
      <v-col cols="12" md="6" style="min-height: 250px">
        <SiteProximityMap :model-value="site?.coordinates ?? {}"></SiteProximityMap>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import SiteAutocomplete from './SiteAutocomplete.vue'
import { SiteInput, SiteItem } from '@/api'
import SitePreviewCard from './SitePreviewCard.vue'
import { useToggle } from '@vueuse/core'
import SiteProximityMap from './SiteProximityMap.vue'

const site = defineModel<SiteItem | SiteInput>()

defineProps<{ class?: any }>()

const [showEdit, toggleEdit] = useToggle(false)

function updateSite(s: SiteItem | SiteInput | undefined) {
  site.value = s
  if (!!s) toggleEdit(false)
}
</script>

<style scoped lang="scss"></style>
