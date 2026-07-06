<template>
  <CenteredSpinner v-if="isPending" text="Loading instance settings..." />
  <v-alert v-else-if="fetchError" color="error" icon="mdi-alert">
    Failed to load instance settings
  </v-alert>
  <v-container v-else-if="instance">
    <v-row>
      <v-col>
        <v-alert v-if="updateError" color="error" icon="mdi-alert">
          Failed to update settings
        </v-alert>
      </v-col>
    </v-row>
    <v-confirm-edit v-model="instance">
      <template #default="{ isPristine, save, cancel, model: proxy, actions }">
        <v-row>
          <v-col cols="12" sm="3" class="px-3 d-flex align-center justify-center">
            <IconEditor />
          </v-col>

          <v-col cols="12" sm="9" variant="text" class="d-flex align-center">
            <div class="w-100">
              <v-text-field
                v-model="proxy.value.title"
                label="Instance name"
                class="pb-4"
                hint="The name that is displayed in the navbar and front page"
                persistent-hint
                v-bind="schema('title')"
              />
              <v-textarea
                v-model="proxy.value.description"
                :rows="2"
                label="Instance description"
                hint="A short description of the database purpose to be displayed on the front page."
                persistent-hint
                v-bind="schema('description')"
                :spellcheck="false"
              />
            </div>
          </v-col>
        </v-row>
        <v-expand-transition>
          <div class="w-100 d-flex justify-end" v-if="!isPristine">
            <v-btn variant="plain" @click="cancel">Cancel</v-btn>
            <v-btn
              text="OK"
              @click="
                mutateAsync(
                  { body: { title: proxy.value.title, description: proxy.value.description } },
                  { onSuccess: save }
                )
              "
            />
          </div>
        </v-expand-transition>
      </template>
    </v-confirm-edit>
    <v-card>
      <v-list>
        <v-list-item>
          <v-switch
            :model-value="instance.is_public"
            @update:model-value="(v) => togglePublicAccess({ body: !!v })"
            label="Instance is public"
            class="mb-5"
            color="primary"
            hint="A private instance requires user authentication to get access to any data. A public instance allows read-only access to anonymous users."
            persistent-hint
          />
        </v-list-item>
        <v-divider />
        <v-list-item>
          <v-switch
            :model-value="instance.account_requests_enabled"
            @update:model-value="(v) => togglePublicRegistration({ body: !!v })"
            label="Allow contributor registration"
            color="primary"
            hint="If enabled, visitors may apply for an account with Contributor privileges."
            persistent-hint
          />
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { $InstanceSettingsUpdate, InstanceSettings } from '@/api'
import {
  getInstanceSettingsQueryKey,
  togglePublicAccessMutation,
  togglePublicRegistrationMutation,
  updateInstanceSettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useSchema } from '@/composables/schema'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useInstanceSettings } from '.'
import IconEditor from './InstanceIcon.vue'

const { instance, reload, isPending, error: fetchError } = useInstanceSettings()

const {
  bind: { schema, field },
  dispatchErrors
} = useSchema($InstanceSettingsUpdate)

const { feedback } = useFeedback()

const queryClient = useQueryClient()

const {
  mutateAsync,
  error: updateError,
  isPending: isUpdating
} = useMutation({
  ...updateInstanceSettingsMutation(),
  onSuccess: () => {
    // queryClient.invalidateQueries({ queryKey: getInstanceSettingsQueryKey() })
    reload()
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors
})

const {
  mutateAsync: togglePublicAccess,
  error: togglePublicAccessError,
  isPending: isTogglingPublicAccess
} = useMutation({
  ...togglePublicAccessMutation(),
  onSuccess: () => {
    // queryClient.invalidateQueries({ queryKey: getInstanceSettingsQueryKey() })
    reload()
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors
})

const {
  mutateAsync: togglePublicRegistration,
  error: togglePublicRegistrationError,
  isPending: isTogglingPublicRegistration
} = useMutation({
  ...togglePublicRegistrationMutation(),
  onSuccess: () => {
    // queryClient.invalidateQueries({ queryKey: getInstanceSettingsQueryKey() })
    reload()
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors
})

// async function setPublic(value: boolean | null) {
//   if (value === null) return
//   await mutateAsync({ body: { public: value } })
// }

async function submit(model: InstanceSettings) {
  await mutateAsync({ body: model })
}
</script>

<style scoped></style>
