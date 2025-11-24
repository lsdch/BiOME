<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} bio material`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>

    <v-container class="bg-main overflow-y-auto responsive-container" fluid min-height="100%">
      <v-row align="stretch" class="bg-main">
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
          <v-card title="Bio material" prepend-icon="mdi-package-variant">
            <v-divider />
            <OccurrenceForm v-model="model.biomaterial.external" />
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

const dialog = defineModel<boolean>('dialog')

const model = defineModel<OccurrenceModel.OccurrenceModel>({
  default: OccurrenceModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: OccurrenceModel.OccurrenceModel | undefined]
}>()
</script>

<style scoped lang="scss"></style>
