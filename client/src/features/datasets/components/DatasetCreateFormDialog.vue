<template>
  <FormDialog v-model="dialog" title="Create dataset" @submit="submit()">
    <template #activator="activatorProps">
      <slot name="activator" v-bind="activatorProps" />
    </template>
    <v-card-text class="d-flex flex-column ga-1 py-1">
      <v-text-field
        label="Label"
        v-model="model.label"
        v-bind="schema('label')"
        :error-messages="errors('label')"
      />
      <v-textarea
        label="Description"
        v-model="model.description"
        v-bind="schema('description')"
        :error-messages="errors('description')"
      />
      <v-switch label="Public" color="primary" v-model="model.public" v-bind="schema('public')" />
      <v-switch label="Pinned" color="primary" v-model="model.pinned" v-bind="schema('pinned')" />
      <!-- <UserPicker
        v-model="model.maintainers"
        label="Maintainers"
        item-value="id"
        v-bind="schema('maintainers')"
      /> -->
    </v-card-text>
  </FormDialog>
</template>

<script setup lang="ts">
import { $DatasetInput, AppError, DatasetInput, DatasetWithMaintainers } from '@/api'
import { createDatasetMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { useSchemaBinding, useSchemaErrors } from '@/composables/schema'
import { useMutation } from '@tanstack/vue-query'
import { ref } from 'vue'

const emit = defineEmits<{
  created: [created: DatasetWithMaintainers]
  error: [error: AppError]
}>()

const dialog = defineModel<boolean>('dialog', { default: false })

const model = ref<DatasetInput>({
  label: '',
  pinned: false,
  maintainers: [],
  public: false
})

const { schema } = useSchemaBinding($DatasetInput)

const { mutateAsync: createDataset, error } = useMutation(createDatasetMutation())
const { errors } = useSchemaErrors($DatasetInput, error)

async function submit() {
  await createDataset(
    { body: model.value },
    {
      onSuccess: (created) => {
        emit('created', created)
        dialog.value = false
      },
      onError: (error) => {
        emit('error', error)
      }
    }
  )
}
</script>

<style scoped lang="scss"></style>
