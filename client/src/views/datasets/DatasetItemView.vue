<template>
  <v-container v-if="error">
    <v-row>
      <v-col>
        {{ error }}
      </v-col>
    </v-row>
  </v-container>
  <v-container v-else-if="dataset" class="fill-height bg-surface align-stretch flex-column" fluid>
    <v-row :class="['justify-start align-start flex-wrap flex-grow-0', { 'fill-height': lgAndUp }]">
      <v-col v-if="editing" cols="12" lg="6" class="align-self-start">
        <DatasetEditForm
          v-model="dataset"
          @cancel="toggleEdit(false)"
          @updated="toggleEdit(false)"
        />
      </v-col>
      <v-col v-else cols="12" lg="6" class="align-self-start">
        <div class="text-h5 d-flex justify-space-between align-center">
          {{ dataset.label }}
          <v-btn
            v-if="isUserMaintainer"
            color="primary"
            icon="mdi-pencil"
            variant="plain"
            @click="toggleEdit(true)"
          />
        </div>
        <v-divider class="my-3" />

        <v-list>
          <v-list-item
            title="Description"
            :subtitle="dataset.description || 'No description'"
            :class="{ empty: !dataset.description }"
          />
          <v-list-item title="Maintainers">
            <template #subtitle>
              <v-chip v-for="(maintainer, key) in dataset.maintainers" :key>
                {{ maintainer.full_name }}
              </v-chip>
            </template>
          </v-list-item>
        </v-list>
        <v-divider class="my-3"></v-divider>
        <div>
          <v-icon class="mx-2">mdi-map-marker</v-icon>
          <span class="text-overline"> {{ dataset.sites?.length }} sites </span>
        </div>
        <div>
          <ItemDateChip v-if="dataset.meta?.created" icon="created" :date="dataset.meta.created" />
          <ItemDateChip
            v-if="dataset.meta?.modified"
            icon="updated"
            :date="dataset.meta.modified"
          />
        </div>
      </v-col>
      <v-col cols="12" lg="6" class="align-self-stretch flex-grow-1 w-100">
        <ResponsiveDialog v-model:open="mobileMap" :as-dialog="!lgAndUp">
          <template #="{ isDialog }">
            <SitesMap
              :items="dataset.sites ?? undefined"
              :closable="isDialog"
              @close="toggleMobileMap(false)"
            >
              <template #marker="{ item }">
                <SitePopup :item />
              </template>
            </SitesMap>
          </template>
        </ResponsiveDialog>
      </v-col>
    </v-row>

    <v-bottom-navigation :active="!lgAndUp">
      <v-btn color="primary" prepend-icon="mdi-map" @click="toggleMobileMap(true)" text="Map" />
    </v-bottom-navigation>
  </v-container>
</template>

<script setup lang="ts">
import { LocationService } from '@/api'
import SitesMap from '@/components/maps/SitesMap.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import ItemDateChip from '@/components/toolkit/ItemDateChip.vue'
import { useUserStore } from '@/stores/user'
import { useToggle } from '@vueuse/core'
import { computed, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'
import DatasetEditForm from './DatasetEditForm.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'

const { user } = useUserStore()

const [editing, toggleEdit] = useToggle(false)

const [mobileMap, toggleMobileMap] = useToggle(false)

const { lgAndUp } = useDisplay()

const { params } = useRoute()
const slug = params.slug as string

const { data: dataset, error } = toRefs(reactive(await fetch()))

async function fetch() {
  return await LocationService.getSiteDataset({ path: { slug } })
}

const isUserMaintainer = computed(() => {
  return !!dataset.value?.maintainers?.find(({ id }) => user?.identity.id === id)
})
</script>

<style lang="scss">
.v-list-item.empty .v-list-item-subtitle {
  font-style: italic;
}
</style>
