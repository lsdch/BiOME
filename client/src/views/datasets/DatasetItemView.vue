<template>
  <div class="d-flex fill-height bg-main overflow-y-auto">
    <v-sheet
      height="fit-content"
      min-height="100%"
      :width="$vuetify.display.lgAndUp ? '50%' : '100%'"
      :class="['d-flex bg-transparent', { 'pa-3': $vuetify.display.lgAndUp }]"
    >
      <PageErrors v-if="error" :error />
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
        :title="baseDataset?.label ?? slug"
        :flat="!$vuetify.display.lgAndUp"
        class="align-self-stretch w-100 d-flex flex-column"
      >
        <template #subtitle>
          <v-chip
            label
            text="Occurrences dataset"
            size="small"
            prepend-icon="mdi-map-marker-multiple"
          />
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
            v-if="baseDataset && userStore.isGranted('Admin')"
            :model-value="baseDataset"
            @update:model-value="({ pinned }) => togglePin(pinned)"
          />
        </template>
        <template #actions>
          <v-spacer />
          <MetaChip v-if="baseDataset" :meta="baseDataset.meta" />
        </template>

        <div class="flex-grow-1">
          <v-img
            :src="`/api/v1/assets/images/datasets/${slug}/${slug}.jpg`"
            :min-height="20"
            :max-height="200"
            cover
            ref="image"
            @error="noImage()"
            gradient="to top, rgba(var(--v-theme-surface)), #00000000"
          >
            <template #error>
              <v-divider class="my-3" />
            </template>
          </v-img>
          <v-card-text v-if="baseDataset?.description" class="text-caption">
            {{ baseDataset?.description }}
          </v-card-text>

          <v-divider />

          <v-list>
            <v-list-item>
              <template #append>
                <span class="text-caption text-muted">Maintainers</span>
              </template>
              <PersonChip
                v-for="(maintainer, key) in baseDataset?.maintainers"
                :person="maintainer"
                class="ma-1"
                :key
              />
            </v-list-item>
          </v-list>

          <v-divider class="mb-3" />
          <slot name="details" />

          <v-divider />
        </div>
      </v-card>
    </v-sheet>

    <ResponsiveDialog v-model:open="mobileMap" :as-dialog="!$vuetify.display.lgAndUp">
      <template #="{ isDialog }">
        <v-sheet
          width="50%"
          :class="['fill-height position-sticky top-0 bg-transparent', { 'py-4 px-3': !isDialog }]"
          max-height="100vh"
        >
          <v-card :rounded="!mobileMap" class="fill-height">
            <slot name="map" :baseDataset :toggleMobileMap :isDialog />
          </v-card>
        </v-sheet>
      </template>
    </ResponsiveDialog>
  </div>
  <v-bottom-navigation :active="!$vuetify.display.lgAndUp">
    <v-btn color="primary" prepend-icon="mdi-map" @click="toggleMobileMap(true)" text="Map" />
  </v-bottom-navigation>
</template>

<script setup lang="ts" generic="DatasetType extends OccurrenceDataset | SiteDataset">
import { Dataset, OccurrenceDataset, SiteDataset } from '@/api'
import { getDatasetOptions } from '@/api/gen/@tanstack/vue-query.gen'
import PersonChip from '@/components/people/PersonChip'
import MetaChip from '@/components/toolkit/MetaChip'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'
import { useUserStore } from '@/stores/user'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { storeToRefs } from 'pinia'
import { computed, useTemplateRef } from 'vue'
import DatasetEditForm from './DatasetEditForm.vue'
import DatasetPinButton from './DatasetPinButton.vue'

const image = useTemplateRef<HTMLImageElement>('image')

function noImage() {
  image.value?.remove()
}

const { slug } = defineProps<{
  slug: string
}>()

const dataset = defineModel<DatasetType | undefined>('dataset')

const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const [editing, toggleEdit] = useToggle(false)

const [mobileMap, toggleMobileMap] = useToggle(false)

const {
  data: baseDataset,
  error,
  isPending,
  isError,
  refetch
} = useQuery(getDatasetOptions({ path: { slug } }))

const queryClient = useQueryClient()

function togglePin(pinned: boolean) {
  queryClient.setQueryData(
    getDatasetOptions({ path: { slug } }).queryKey,
    (oldData: Dataset | undefined) => {
      return oldData ? { ...oldData, pinned } : undefined
    }
  )
}

const isUserMaintainer = computed(() => {
  return !!baseDataset.value?.maintainers?.find(({ id }) => user.value?.identity.id === id)
})
</script>

<style lang="scss">
.v-list-item.empty .v-list-item-subtitle {
  font-style: italic;
}
</style>
