<template>
  <v-row>
    <v-col>
      <v-text-field label="Code" v-model.trim="model.code"></v-text-field>
    </v-col>
  </v-row>
  <v-card title="Identification" flat>
    <IdentificationFormFields v-model="model.identification!" />
  </v-card>
  <v-card title="Content" flat>
    <v-row>
      <v-col cols="12" md="3">
        <ExtBioMatQuantityPicker v-model="model.quantity" label="Specimen quantity" />
      </v-col>
      <v-col cols="12" md="9">
        <v-text-field v-model="model.content_description" label="Additional details (optional)" />
      </v-col>
    </v-row>
  </v-card>
  <v-row>
    <v-col>
      <v-card title="References " flat>
        <v-text-field v-model="model.collection" label="Collection" />
        <v-combobox v-model="model.vouchers" label="Item vouchers" />
        <ArticlesPicker v-model="model.published_in" label="Published in" multiple clearable />
      </v-card>
    </v-col>
  </v-row>
  <v-divider></v-divider>
  <v-row>
    <v-col>
      <v-textarea label="Comments (optional)" v-model.trim="model.comments"></v-textarea>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import {
  $ExternalBioMatInput,
  $ExternalBioMatUpdate,
  BioMaterialWithDetails,
  DateWithPrecision,
  ExternalBioMatInput,
  ExternalBioMatUpdate,
  IdentificationInput,
  SamplesService
} from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import { reactiveComputed, useToggle } from '@vueuse/core'
import ArticlesPicker from '../references/ArticlesPicker.vue'
import ExtBioMatQuantityPicker from './ExtBioMatQuantityPicker.vue'
import IdentificationFormFields from './IdentificationFormFields.vue'

const props = defineProps<FormProps<BioMaterialWithDetails>>()
const emit = defineEmits<FormEmits<BioMaterialWithDetails>>()

const initial: Partial<Omit<ExternalBioMatInput, 'identification'>> & {
  identification: Partial<IdentificationInput>
} = {
  identification: { identified_on: { date: {}, precision: 'Day' } }
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({
    sampling,
    code,
    comments,
    identification: { identified_by, identified_on, taxon },
    is_type,
    external: { archive, quantity, content_description, original_link, original_taxon },
    published_in
  }: BioMaterialWithDetails & {
    category: 'External'
    external: Exclude<BioMaterialWithDetails['external'], null | undefined>
  }): ExternalBioMatUpdate {
    return {
      sampling_id: sampling.id,
      code,
      is_type,
      comments,
      identification: {
        identified_by: identified_by.alias,
        identified_on: DateWithPrecision.toInput(identified_on),
        taxon: taxon.name
      },
      collection: archive.collection,
      vouchers: archive.vouchers,
      quantity,
      content_description,
      original_link,
      original_taxon,
      published_in: published_in.map(({ code }) => code)
    }
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $ExternalBioMatInput : $ExternalBioMatUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    // TODO: check if we can use form validation for better typing
    create: (body) =>
      SamplesService.createExternalBioMat({ body: body as unknown as ExternalBioMatInput }),
    edit: ({ code }, body) => SamplesService.updateExternalBioMat({ body, path: { code } })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
