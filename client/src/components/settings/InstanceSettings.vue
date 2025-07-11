<template>
  <CenteredSpinner v-if="isPending" text="Loading instance settings..." />
  <v-alert v-else-if="fetchError" color="error" icon="mdi-alert">
    Failed to load instance settings
  </v-alert>
  <v-confirm-edit v-else v-model="instance">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @reset="cancel()"
        @submit="submit(proxy.value).then(() => save())"
        :loading="isUpdating"
      />
      <v-container>
        <v-row>
          <v-col>
            <v-alert v-if="updateError" color="error" icon="mdi-alert">
              Failed to update settings
            </v-alert>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" sm="3" class="px-3 d-flex align-center justify-center">
            <IconEditor />
          </v-col>
          <v-col cols="12" sm="9" variant="text" class="d-flex align-center">
            <div class="w-100">
              <v-text-field
                v-model="proxy.value.name"
                label="Instance name"
                class="pb-4"
                hint="The name that is displayed in the navbar and front page"
                persistent-hint
                v-bind="schema('name')"
              />
              <v-text-field
                v-model="proxy.value.description"
                label="Instance description"
                hint="A short description of the database purpose to be displayed on the front page."
                persistent-hint
                clearable
                v-bind="schema('description')"
              />
            </div>
          </v-col>
        </v-row>
        <v-divider />
        <v-switch
          v-model="proxy.value.public"
          label="Instance is public"
          class="mb-5"
          color="primary"
          hint="A private instance requires user authentication to get access to any data. A public instance allows read-only access to anonymous users on a subset of pages."
          persistent-hint
        />
        <v-divider />
        <v-switch
          v-model="proxy.value.allow_contributor_signup"
          label="Allow contributor registration"
          color="primary"
          hint="If enabled, visitors may apply for an account with Contributor privileges."
          persistent-hint
        />
      </v-container>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts">
import { $InstanceSettingsInput, InstanceSettings } from '@/api'
import {
  instanceSettingsQueryKey,
  updateInstanceSettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useInstanceSettings } from '.'
import { useSchema } from '../../composables/schema'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import IconEditor from './InstanceIcon.vue'
import SettingsFormActions from './SettingsFormActions.vue'

const { instance, reload, isPending, error: fetchError } = useInstanceSettings()

const {
  bind: { schema },
  dispatchErrors
} = useSchema($InstanceSettingsInput)

const { feedback } = useFeedback()

const queryClient = useQueryClient()

const {
  mutateAsync,
  error: updateError,
  isPending: isUpdating
} = useMutation({
  ...updateInstanceSettingsMutation(),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: instanceSettingsQueryKey() })
    reload()
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors
})

async function submit(model: InstanceSettings) {
  await mutateAsync({ body: model })
}
</script>

<style scoped></style>
