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
          @click="(save(), submit())"
          :disabled="isPristine"
          text="Save"
        />
        <v-btn color="" @click="(cancel(), emit('cancel'))" text="Cancel" />
      </div>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts">
import { $DatasetUpdate, SiteDataset, DatasetUpdate, DatasetsService } from '@/api'
import PersonPicker from '@/components/people/PersonPicker.vue'
import { useSchema } from '@/components/toolkit/forms/schema'
import { useFeedback } from '@/stores/feedback'

const dataset = defineModel<SiteDataset>()
const { schema, errorHandler } = useSchema($DatasetUpdate)
const { feedback } = useFeedback()

const emit = defineEmits<{
  cancel: []
  updated: [dataset: SiteDataset]
}>()

async function submit() {
  if (!dataset.value) return
  const { label, description } = dataset.value
  const body: DatasetUpdate = {
    label,
    description,
    maintainers: dataset.value?.maintainers?.map(({ alias }) => alias) || null
  }
  await DatasetsService.updateSiteDataset({ path: { slug: dataset.value.slug }, body })
    .then(errorHandler)
    .then((updated) => {
      dataset.value = updated
      feedback({ type: 'success', message: 'Updated dataset infos' })
      emit('updated', updated)
    })
}
</script>

<style scoped></style>
