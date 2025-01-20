<template>
  <v-confirm-edit v-model="model">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @reset="reloadSettings().then(cancel)"
        @submit="(save(), submit())"
      />
      <v-container>
        <v-row>
          <v-col>
            <v-card title="Geocoding" flat>
              <v-card-text>
                <p class="mb-3">
                  Geoapify provides geocoding and location services. When enabled, some sites
                  metadata may be inferred from their coordinates (unless explicitly specified by
                  users).
                </p>

                <p class="mb-5">
                  Using this service requires an API key, which can be obtained by registering an
                  account at <a href="https://www.geoapify.com/">Geoapify</a>.
                </p>
                <v-text-field
                  v-model.trim="proxy.value.geoapify_api_key"
                  label="API key"
                  :step="1"
                  v-bind="field('geoapify_api_key')"
                  control-variant="stacked"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts">
import { $ServiceSettingsUpdate, SettingsService } from '@/api'
import { handleErrors } from '@/api/responses'
import { useFeedback } from '@/stores/feedback'
import { ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import SettingsFormActions from './SettingsFormActions.vue'

const { feedback } = useFeedback()

async function fetch() {
  return SettingsService.serviceSettings().then(
    handleErrors((err) => {
      feedback({ message: 'Failed to retrieve service settings', type: 'error' })
      console.error(err)
    })
  )
}

async function reloadSettings() {
  model.value = await fetch()
}

async function submit() {
  return SettingsService.updateServiceSettings({ body: model.value })
    .then(errorHandler)
    .then(() => feedback({ message: 'Settings updated', type: 'success' }))
}

const model = ref(await fetch())

const { field, errorHandler } = useSchema($ServiceSettingsUpdate)
</script>

<style lang="scss">
.settings-list {
  .v-list-item-title {
    margin-top: 10px;
    margin-bottom: 20px;
  }
}
.v-slider-track__tick-label {
  font-size: small;
}
</style>
