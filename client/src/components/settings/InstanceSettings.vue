<template>
  <SettingsForm @reset="reloadSettings" @submit="submit">
    <v-row>
      <v-col cols="12" sm="3" class="px-3 d-flex align-center justify-center">
        <IconEditor />
      </v-col>
      <v-col cols="12" sm="9" variant="text" class="pt-8 d-flex flex-column justify-center">
        <v-text-field
          v-model="model.name"
          label="Instance name"
          class="pb-4"
          hint="The name that is displayed in the navbar and front page"
          persistent-hint
          v-bind="field('name')"
        />
        <v-text-field
          v-model="model.description"
          label="Instance description"
          class="mb-5"
          hint="A short description of the database purpose to be displayed on the front page."
          persistent-hint
          clearable
          v-bind="field('description')"
        />
      </v-col>
    </v-row>
    <v-divider />
    <v-switch
      v-model="model.public"
      label="Instance is public"
      class="mb-5"
      color="primary"
      hint="A private instance requires user authentication to get access to any data. A public instance allows read-only access to anonymous users on a subset of pages."
      persistent-hint
    />
    <v-divider />
    <v-switch
      v-model="model.allow_contributor_signup"
      label="Allow contributor registration"
      color="primary"
      hint="If enabled, visitors may apply for an account with Contributor privileges."
      persistent-hint
    />
  </SettingsForm>
</template>

<script setup lang="ts">
import { $InstanceSettingsInput, SettingsService } from '@/api'
import { ref } from 'vue'
import { useInstanceSettings } from '.'
import { useSchema } from '../toolkit/forms/schema'
import IconEditor from './InstanceIcon.vue'
import SettingsForm from './SettingsForm.vue'
import { useFeedback } from '@/stores/feedback'

const { settings, reload } = useInstanceSettings()

const model = ref(settings)

const { field, errorHandler } = useSchema($InstanceSettingsInput)

async function reloadSettings() {
  model.value = await reload()
}

const { feedback } = useFeedback()

async function submit() {
  await SettingsService.updateInstanceSettings({ body: model.value })
    .then(errorHandler)
    .then(() => {
      reloadSettings()
      feedback({ message: 'Updated settings', type: 'success' })
    })
}
</script>

<style scoped></style>
