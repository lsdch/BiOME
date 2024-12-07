<template>
  <v-layout class="fill-height bg-surface">
    <PageErrors v-if="error !== undefined" :error />
    <v-container v-else-if="dataset" class="align-stretch flex-column" fluid>
      <v-row
        :class="['justify-start align-start flex-wrap flex-grow-0', { 'fill-height': lgAndUp }]"
      >
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
          <div>
            <MetaChip :meta="dataset.meta" />
          </div>
          <v-divider class="my-3" />
          <DatasetTabs :dataset />
        </v-col>
        <v-col cols="12" lg="6" class="align-self-stretch flex-grow-1 w-100">
          <ResponsiveDialog v-model:open="mobileMap" :as-dialog="!lgAndUp">
            <template #="{ isDialog }">
              <v-card class="fill-height" :rounded="!mobileMap">
                <SitesMap
                  :items="dataset.sites ?? undefined"
                  :closable="isDialog"
                  @close="toggleMobileMap(false)"
                >
                  <template #marker="{ item }">
                    <SitePopup :item />
                  </template>
                </SitesMap>
              </v-card>
            </template>
          </ResponsiveDialog>
        </v-col>
      </v-row>

      <v-bottom-navigation :active="!lgAndUp">
        <v-btn color="primary" prepend-icon="mdi-map" @click="toggleMobileMap(true)" text="Map" />
      </v-bottom-navigation>
    </v-container>
  </v-layout>
</template>

<script setup lang="ts">
import { DatasetsService } from '@/api'
import SitesMap from '@/components/maps/SitesMap.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import ItemDateChip from '@/components/toolkit/ItemDateChip.vue'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'
import { useUserStore } from '@/stores/user'
import { useToggle } from '@vueuse/core'
import { computed, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'
import DatasetEditForm from './DatasetEditForm.vue'
import DatasetTabs from './DatasetTabs.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'

const { user } = useUserStore()

const [editing, toggleEdit] = useToggle(false)

const [mobileMap, toggleMobileMap] = useToggle(false)

const { lgAndUp } = useDisplay()

const { params } = useRoute()
const slug = params.slug as string

const { data: dataset, error } = toRefs(reactive(await fetch()))

async function fetch() {
  return await DatasetsService.getDataset({ path: { slug } })
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
