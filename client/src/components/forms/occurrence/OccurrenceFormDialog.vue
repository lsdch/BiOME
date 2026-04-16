<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} bio material`"
    :max-width="1200"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-container class="bg-main overflow-y-auto responsive-container" fluid min-height="100%">
      <v-row>
        <v-col>
          <v-card
            title="Datasets"
            subtitle="Optionally add the new occurrence to one or more datasets"
            class="small-card-title"
          >
            <template #append>
              <v-switch v-model="withDatasets" color="primary" hide-details></v-switch>
            </template>
            <template v-if="withDatasets">
              <v-divider></v-divider>
              <v-card-text>
                <DatasetPicker label="" multiple clearable>
                  <template #append>
                    <v-btn text="New" prepend-icon="mdi-plus" variant="tonal" rounded="md"></v-btn>
                  </template>
                </DatasetPicker>
              </v-card-text>
            </template>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="bg-main align-stretch">
        <v-col cols="12" md="6">
          <SiteFormComponent class="fill-height small-card-title" v-model="model.site" />
        </v-col>
        <v-col cols="12" md="6">
          <div class="d-flex flex-column ga-3">
            <SamplingFormComponent :site="model.site" v-model="model.sampling" />
          </div>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card title="Occurrence identification" class="" prepend-icon="mdi-package-variant">
            <v-divider />
            <OccurrenceForm v-model="model.biomaterial" />
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { FormProps } from '@/lib/mutations'
import { OccurrenceModel } from '@/models'
import OccurrenceForm from './OccurrenceForm.vue'
import SamplingFormComponent from './OccurrenceFormSampling.vue'
import SiteFormComponent from './OccurrenceFormSite.vue'
import DatasetPicker from '@/features/datasets/components/DatasetPicker.vue'
import { ref } from 'vue'

const dialog = defineModel<boolean>('dialog')

const model = defineModel<OccurrenceModel.OccurrenceModel>({
  default: OccurrenceModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: OccurrenceModel.OccurrenceModel | undefined]
}>()

const withDatasets = ref(false)
</script>

<style scoped lang="scss"></style>
