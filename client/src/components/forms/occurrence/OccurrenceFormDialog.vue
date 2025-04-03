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

    <v-container class="bg-main overflow-y-auto" fluid min-height="100%">
      <v-row align="stretch" class="bg-main">
        <v-col cols="12" md="6">
          <SiteFormComponent class="fill-height small-card-title" v-model="model.site" />
        </v-col>
        <v-col cols="12" md="6">
          <div class="d-flex flex-column ga-3">
            <EventFormComponent :site="model.site" v-model="model.event" />
            <SamplingFormComponent :event="model.event" v-model="model.sampling" />
          </div>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card title="Bio material" prepend-icon="mdi-package-variant">
            <template #append v-if="mode === 'Create' && $vuetify.display.smAndUp">
              <OccurrenceCategoryBtnToggle v-model="model.biomaterial.category" />
            </template>
            <v-divider v-if="model.biomaterial.category" class="mb-3" />
            <ExternalBioMatForm
              v-if="model.biomaterial.category === 'External'"
              v-model="model.biomaterial.external"
            />
            <v-card-text v-else-if="model.biomaterial.category === 'Internal'">
              [Internal bio mat form]
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { FormProps } from '@/functions/mutations'
import { OccurrenceModel } from '@/models'
import ExternalBioMatForm from './ExternalBioMatForm.vue'
import { OccurrenceCategoryBtnToggle } from './OccurrenceCategoryBtnToggle'
import EventFormComponent from './OccurrenceFormEvent.vue'
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
