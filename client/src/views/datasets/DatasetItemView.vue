<template>
  <div class="d-flex fill-height bg-main overflow-y-auto">
    <v-sheet
      height="fit-content"
      min-height="100%"
      :width="lgAndUp ? '50%' : '100%'"
      :class="['d-flex bg-transparent', { 'pa-3': lgAndUp }]"
    >
      <PageErrors v-if="isError && error" :error />
      <v-card v-else-if="editing && dataset" class="align-self-stretch w-100 d-flex flex-column">
        <v-card-text>
          <DatasetEditForm
            v-model="dataset"
            @cancel="toggleEdit(false)"
            @updated="toggleEdit(false)"
          />
        </v-card-text>
      </v-card>

      <v-card
        v-else
        min-height="100%"
        :title="dataset?.label ?? slug"
        :flat="!lgAndUp"
        class="align-self-stretch w-100 d-flex flex-column"
      >
        <template #subtitle>
          <v-chip label text="Sites dataset" size="small" prepend-icon="mdi-map-marker-multiple" />
        </template>
        <template #prepend>
          <v-avatar variant="outlined">
            <v-icon icon="mdi-folder-table" @click="refetch()" />
          </v-avatar>
        </template>
        <template #append>
          <v-btn
            v-if="isUserMaintainer || userStore.isGranted('Admin')"
            color="primary"
            icon="mdi-pencil"
            variant="plain"
            @click="toggleEdit(true)"
          />
          <DatasetPinButton
            v-if="dataset && userStore.isGranted('Admin')"
            :model-value="dataset"
            @update:model-value="({ pinned }) => togglePin(pinned)"
          />
        </template>
        <template #actions>
          <v-spacer />
          <MetaChip v-if="dataset" :meta="dataset.meta" />
        </template>

        <v-divider class="my-3" />

        <CenteredSpinner v-if="isPending" :height="500" size="large" color="primary" />
        <div v-else-if="dataset" class="flex-grow-1">
          <v-card-text v-if="dataset.description" class="text-caption font-weight-thin">
            {{ dataset.description }}
          </v-card-text>

          <v-list>
            <v-list-item title="Maintainers">
              <template #subtitle>
                <PersonChip
                  v-for="(maintainer, key) in dataset.maintainers"
                  :person="maintainer"
                  class="ma-1"
                  :key
                />
              </template>
            </v-list-item>
          </v-list>

          <v-divider class="my-3" />
          <div class="flex-grow-1">
            <DatasetTabs :dataset flat />
          </div>
          <v-divider />
        </div>
      </v-card>
    </v-sheet>

    <ResponsiveDialog v-model:open="mobileMap" :as-dialog="!lgAndUp">
      <template #="{ isDialog }">
        <v-sheet
          width="50%"
          :class="['fill-height position-sticky top-0 bg-transparent', { 'py-4 px-3': !isDialog }]"
          max-height="100vh"
        >
          <v-card :rounded="!mobileMap" class="fill-height">
            <SitesMap
              :items="dataset?.sites"
              :closable="isDialog"
              @close="toggleMobileMap(false)"
              clustered
            >
              <template #popup="{ item }">
                <SitePopup :item />
              </template>
            </SitesMap>
          </v-card>
        </v-sheet>
      </template>
    </ResponsiveDialog>
  </div>
  <v-bottom-navigation :active="!lgAndUp">
    <v-btn color="primary" prepend-icon="mdi-map" @click="toggleMobileMap(true)" text="Map" />
  </v-bottom-navigation>
</template>

<script setup lang="ts" generic="DatasetType extends OccurrenceDataset | SiteDataset">
import { ErrorModel, OccurrenceDataset, SiteDataset } from '@/api'
import SitesMap from '@/components/maps/SitesMap.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'
import { useUserStore } from '@/stores/user'
import { Options, OptionsLegacyParser } from '@hey-api/client-fetch'
import { UndefinedInitialQueryOptions, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useDisplay } from 'vuetify'
import DatasetPinButton from './DatasetPinButton.vue'
import DatasetTabs from './DatasetTabs.vue'
import DatasetEditForm from './DatasetEditForm.vue'

interface DatasetQueryData {
  headers?: {
    Authorization: string
  }
  path: {
    slug: string
  }
}

type QueryKey<TOptions extends Options> = [
  Pick<TOptions, 'baseUrl' | 'body' | 'headers' | 'path' | 'query'> & {
    _id: string
    _infinite?: boolean
  }
]

const { slug, query } = defineProps<{
  slug: string
  query: (options: OptionsLegacyParser<DatasetQueryData>) => UndefinedInitialQueryOptions<
    DatasetType,
    ErrorModel,
    DatasetType,
    QueryKey<DatasetQueryData>
  > & {
    queryKey: QueryKey<DatasetQueryData>
  }
}>()

const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const [editing, toggleEdit] = useToggle(false)

const [mobileMap, toggleMobileMap] = useToggle(false)

const { lgAndUp } = useDisplay()

const { data: dataset, error, isPending, isError, refetch } = useQuery(query({ path: { slug } }))

const queryClient = useQueryClient()

function togglePin(pinned: boolean) {
  queryClient.setQueryData(query({ path: { slug } }).queryKey, (oldData: DatasetType) => {
    return { ...oldData, pinned }
  })
}

const isUserMaintainer = computed(() => {
  return !!dataset.value?.maintainers?.find(({ id }) => user.value?.identity.id === id)
})
</script>

<style lang="scss">
.v-list-item.empty .v-list-item-subtitle {
  font-style: italic;
}
</style>
