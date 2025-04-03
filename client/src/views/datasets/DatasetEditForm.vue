<template>
  <v-confirm-edit v-model="dataset">
    <template v-slot:default="{ model: proxy, save, cancel, isPristine, actions: _ }">
      <v-text-field v-model="proxy.value.label" label="Label" v-bind="schema('label')" />
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
        restrict="Contributor"
        return-objects
        v-bind="schema('maintainers')"
        clearable
      />
      <div class="d-flex justify-end">
        <v-btn
          color="primary"
          class="mx-3"
          @click="submit(proxy.value).then(() => save())"
          :disabled="isPristine"
          text="Save"
          :loading
        />
        <v-btn color="" @click="(cancel(), emit('cancel'))" text="Cancel" />
      </div>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts" generic="DatasetType extends OccurrenceDataset | SiteDataset">
import { $DatasetUpdate, DatasetUpdate, OccurrenceDataset, SiteDataset } from '@/api'
import { updateDatasetMutation } from '@/api/gen/@tanstack/vue-query.gen'
import PersonPicker from '@/components/people/PersonPicker.vue'
import { useSchema } from '@/composables/schema'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'

const dataset = defineModel<DatasetType>({ required: true })
const {
  bind: { schema },
  dispatchErrors
} = useSchema($DatasetUpdate)
const { feedback } = useFeedback()

const emit = defineEmits<{
  cancel: []
  updated: [dataset: DatasetType]
}>()

const {
  mutateAsync,
  isPending: loading,
  error
} = useMutation({
  ...updateDatasetMutation(),
  onSuccess(updated) {
    dataset.value = { ...dataset.value, ...updated }
    emit('updated', dataset.value)
    feedback({ message: 'Dataset updated', type: 'success' })
  },
  onError(error) {
    dispatchErrors(error)
  }
})

async function submit({ label, description, maintainers }: DatasetType) {
  const { slug } = dataset.value
  const body: DatasetUpdate = {
    label,
    description,
    maintainers: maintainers?.map(({ alias }) => alias) || null
  }
  return mutateAsync({ path: { slug }, body })
}
</script>

<style scoped></style>
