<template>
  <v-container v-if="error">
    <v-row>
      <v-col>
        {{ error }}
      </v-col>
    </v-row>
  </v-container>
  <v-container v-else-if="dataset" class="fill-height bg-surface align-stretch flex-column" fluid>
    <v-confirm-edit v-model="dataset">
      <template v-slot:default="{ model: proxy, save, cancel, isPristine, actions: _ }">
        <v-row
          :class="['justify-start align-start flex-wrap flex-grow-0', { 'fill-height': lgAndUp }]"
        >
          <v-col cols="12" lg="6" class="align-self-start">
            <v-text-field
              v-if="editing"
              v-model="proxy.value.label"
              label="Label"
              v-bind="schema('label')"
            />
            <div v-else>
              <div class="text-h5 d-flex justify-space-between align-center">
                {{ dataset.label }}
                <v-btn
                  color="primary"
                  icon="mdi-pencil"
                  variant="plain"
                  @click="toggleEdit(true)"
                />
              </div>
              <v-divider class="my-3" />
            </div>

            <!-- EDIT -->
            <div v-if="editing">
              <v-textarea
                v-model="proxy.value.description"
                label="Description"
                variant="outlined"
                v-bind="schema('description')"
              />
              <PersonPicker
                label="Maintainers"
                v-model="proxy.value.maintainers"
                multiple
                users="Contributor"
                return-objects
                v-bind="schema('maintainers')"
                clearable
              />
              <div class="d-flex justify-end">
                <v-btn
                  color="primary"
                  class="mx-3"
                  @click="
                    () => {
                      save()
                      submit()
                    }
                  "
                  :disabled="isPristine"
                  text="Save"
                />
                <v-btn
                  color=""
                  @click="
                    () => {
                      cancel()
                      toggleEdit(false)
                    }
                  "
                  text="Cancel"
                />
              </div>
            </div>

            <!-- SHOW  -->
            <v-list v-else>
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
              <ItemDateChip
                v-if="dataset.meta?.created"
                icon="created"
                :date="dataset.meta.created"
              />
              <ItemDateChip
                v-if="dataset.meta?.modified"
                icon="updated"
                :date="dataset.meta.modified"
              />
            </div>
          </v-col>
          <v-col cols="12" lg="6" class="align-self-stretch flex-grow-1 w-100">
            <SitesMap :items="dataset.sites ?? undefined" v-if="lgAndUp">
              <template #marker="{ item }">
                <SitePopup :item />
              </template>
            </SitesMap>
          </v-col>
        </v-row>
      </template>
    </v-confirm-edit>

    <v-row>
      <v-col>
        <SitesMap v-show="!(lgAndUp || mobile)" :items="dataset.sites ?? undefined">
          <template #marker="{ item }">
            <SitePopup :item />
          </template>
        </SitesMap>
      </v-col>
    </v-row>
    <v-bottom-navigation :active="mobile">
      <v-dialog v-model="mobileMap" fullscreen>
        <SitesMap :items="dataset.sites ?? undefined" closable @close="toggleMobileMap(false)">
          <template #marker="{ item }">
            <SitePopup :item />
          </template>
        </SitesMap>
        <template #activator>
          <v-btn color="primary" icon="mdi-map" @click="toggleMobileMap(true)" />
        </template>
      </v-dialog>
    </v-bottom-navigation>
  </v-container>
</template>

<script setup lang="ts">
import { $SiteDatasetUpdate, LocationService, SiteDatasetUpdate } from '@/api'
import SitesMap from '@/components/maps/SitesMap.vue'
import PersonPicker from '@/components/people/PersonPicker.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import { useSchema } from '@/components/toolkit/forms/form'
import ItemDateChip from '@/components/toolkit/ItemDateChip.vue'
import { useToggle } from '@vueuse/core'
import { computed, reactive, toRefs } from 'vue'
import { useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'

const [editing, toggleEdit] = useToggle(false)

const [mobileMap, toggleMobileMap] = useToggle(false)

const { smAndDown, mdAndDown, lgAndUp, mobile } = useDisplay()

const { params } = useRoute()
const slug = params.slug as string

const fullscreenActive = computed(() => !!document.fullscreenElement)

const { data: dataset, error } = toRefs(reactive(await fetch()))

async function fetch() {
  return await LocationService.getSiteDataset({ path: { slug } })
}

const { schema, errorHandler } = useSchema($SiteDatasetUpdate)

async function submit() {
  if (!dataset.value) return
  const { label, description } = dataset.value
  const body: SiteDatasetUpdate = {
    label,
    description,
    maintainers: dataset.value?.maintainers?.map(({ alias }) => alias) || null
  }
  await LocationService.updateSiteDataset({ path: { slug }, body })
    .then(errorHandler)
    .then((updated) => (dataset.value = updated))
}
</script>

<style lang="scss">
.v-list-item.empty .v-list-item-subtitle {
  font-style: italic;
}
</style>
