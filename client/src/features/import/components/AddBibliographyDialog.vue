<template>
  <CardDialog v-bind="$attrs" title="Add bibliography" :width="800">
    <v-card-text>
      <div class="d-flex ga-3 py-3">
        <v-file-input
          label="Occurrences TSV file"
          v-model="model.file"
          class="required"
          :rules="[(v) => !!v || 'File is required']"
          filter-by-type="text/tab-separated-values, text/csv"
        ></v-file-input>
        <CSVQuotePicker v-model="model.quotes" class="flex-grow-0"></CSVQuotePicker>
      </div>
    </v-card-text>
    <template #append>
      <v-btn
        text="Submit"
        prepend-icon="mdi-upload"
        variant="elevated"
        :loading="isSubmitting"
        @click="submit()"
        :disabled="!model.file"
      ></v-btn>
    </template>
    <template #activator="props">
      <slot name="activator" v-bind="props"></slot>
    </template>
  </CardDialog>
</template>

<script setup lang="ts">
import { addBibliographyCsvMutation } from '@/api/gen/@tanstack/vue-query.gen'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import CSVQuotePicker from '@/components/toolkit/ui/exports/CSVQuotePicker.vue'
import { useMutation } from '@tanstack/vue-query'
import { ref } from 'vue'

const props = defineProps<{ import_id: UUID } & CardDialogProps>()

const model = ref<{ file: File | null; quotes: string }>({ file: null, quotes: '"' })

async function submit() {
  if (!model.value.file) {
    throw new Error('File is required')
  }
  await submitBibliography({
    path: { id: props.import_id },
    body: { file: model.value.file, quotes: model.value.quotes, separator: '\t' }
  })
}

const { mutateAsync: submitBibliography, isPending: isSubmitting } = useMutation(
  addBibliographyCsvMutation()
)
</script>

<style scoped lang="scss"></style>
