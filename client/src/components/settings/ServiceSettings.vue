<template>
  <CenteredSpinner v-if="isPending" text="Loading services settings..." />
  <v-alert v-else-if="error" color="error" icon="mdi-alert">
    Failed to load services settings
  </v-alert>
  <v-confirm-edit v-else v-model="model">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @reset="cancel()"
        @submit="submit(proxy.value).then(() => save())"
        :loading="isUpdating"
      />
      <v-row>
        <v-col>
          <v-alert v-if="updateError" color="error" icon="mdi-alert">
            Failed to update settings
          </v-alert>
        </v-col>
      </v-row>
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
import { $ServiceSettingsUpdate, ServiceSettingsUpdate } from '@/api'
import {
  serviceSettingsOptions,
  updateServiceSettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { useSchema } from '../toolkit/forms/schema'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import SettingsFormActions from './SettingsFormActions.vue'

const { feedback } = useFeedback()

const { data: model, error, refetch, isPending } = useQuery(serviceSettingsOptions())

const { field, dispatchErrors } = useSchema($ServiceSettingsUpdate)

const {
  mutateAsync,
  error: updateError,
  isPending: isUpdating
} = useMutation({
  ...updateServiceSettingsMutation(),
  onSuccess: () => {
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors
})

async function submit(model: ServiceSettingsUpdate) {
  await mutateAsync({ body: model })
}
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
