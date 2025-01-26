<template>
  <v-layout class="fill-height bg-surface overflow-y-auto">
    <PageErrors v-if="error !== undefined" :error />
    <v-container
      v-else-if="dataset"
      class="align-stretch flex-column"
      fluid
      :width="lgAndUp ? '50%' : '100%'"
    >
      <v-row class="justify-start align-start flex-wrap flex-grow-0 overflow-y-scroll">
        <v-col cols="12">
          <v-card v-if="editing" cols="12" lg="6" class="align-self-start">
            <!-- <DatasetEditForm
              v-model="dataset"
              @cancel="toggleEdit(false)"
              @updated="toggleEdit(false)"
            /> -->
          </v-card>
          <v-card
            v-else
            cols="12"
            lg="6"
            :title="dataset.label"
            class="align-self-start d-flex flex-column"
          >
            <template #subtitle>
              <v-chip
                label
                text="Site dataset"
                size="small"
                prepend-icon="mdi-map-marker-multiple"
              />
            </template>
            <template #prepend>
              <v-avatar variant="outlined">
                <v-icon icon="mdi-folder-table" />
              </v-avatar>
            </template>
            <template #append>
              <v-btn
                v-if="isUserMaintainer"
                color="primary"
                icon="mdi-pencil"
                variant="plain"
                @click="toggleEdit(true)"
              />
            </template>
            <template #actions>
              <v-spacer />
              <MetaChip :meta="dataset.meta" />
            </template>

            <v-divider class="my-3" />

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
          </v-card>
        </v-col>
        <!-- <v-col cols="12" lg="6" class="flex-grow-1 position-sticky top-0" style="height: 90vh"> -->

        <!-- </v-col> -->
      </v-row>
    </v-container>
    <ResponsiveDialog v-model:open="mobileMap" :as-dialog="!lgAndUp">
      <template #="{ isDialog }">
        <v-sheet
          width="50%"
          :class="['fill-height position-sticky top-0', { 'pa-3': !isDialog }]"
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
  </v-layout>
  <v-bottom-navigation :active="!lgAndUp">
    <v-btn color="primary" prepend-icon="mdi-map" @click="toggleMobileMap(true)" text="Map" />
  </v-bottom-navigation>
</template>

<script setup lang="ts">
import { DatasetsService } from '@/api'
import SitesMap from '@/components/maps/SitesMap.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'
import { useUserStore } from '@/stores/user'
import { useToggle } from '@vueuse/core'
import { computed, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'
import DatasetTabs from './DatasetTabs.vue'

const { user } = useUserStore()

const [editing, toggleEdit] = useToggle(false)

const [mobileMap, toggleMobileMap] = useToggle(false)

const { lgAndUp } = useDisplay()

const { params } = useRoute()
const slug = params.slug as string

const { data: dataset, error } = toRefs(reactive(await fetch()))

async function fetch() {
  return await DatasetsService.getSiteDataset({ path: { slug } })
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
