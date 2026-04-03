<template>
  <v-list>
    <v-list-item
      title="Update occurrence codes"
      subtitle="This will update occurrence codes based on the current taxon and sampling data."
    >
      <template #append>
        <v-btn
          text="Update codes"
          variant="tonal"
          rounded="md"
          prepend-icon="mdi-refresh"
          :loading="isPending"
          @click="onUpdateCodes"
        />
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { OccurrenceDataset } from '@/api'
import { updateOccurrenceCodesInDatasetMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation } from '@tanstack/vue-query'

const { dataset } = defineProps<{
  dataset: OccurrenceDataset
}>()

const { mutate: updateOccurrenceCodesInDataset, isPending } = useMutation(
  updateOccurrenceCodesInDatasetMutation()
)

const onUpdateCodes = () => {
  if (!dataset.slug) return

  updateOccurrenceCodesInDataset(
    {
      path: {
        slug: dataset.slug
      }
    },
    {
      onSuccess: (changes) => {
        console.info(`${changes.length} occurrence codes updated successfully`)
        emit('refresh')
      }
    }
  )
}

const emit = defineEmits<{
  refresh: []
}>()
</script>

<style scoped lang="scss"></style>
