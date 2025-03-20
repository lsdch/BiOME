<template>
  <CreateUpdateForm
    v-model="item"
    :initial
    :update-transformer
    :create
    :update
    @success="dialog = false"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <v-row>
        <v-col>
          <v-text-field label="Code" v-model.trim="model.code"></v-text-field>
        </v-col>
      </v-row>
      <v-card title="Identification" class="small-card-title" flat>
        <IdentificationFormFields v-model="model.identification!" />
      </v-card>
      <v-card title="Content" class="small-card-title" flat>
        <v-row>
          <v-col cols="12" md="3">
            <ExtBioMatQuantityPicker v-model="model.quantity" label="Specimen quantity" />
          </v-col>
          <v-col cols="12" md="9">
            <v-text-field
              v-model="model.content_description"
              label="Additional details (optional)"
            />
          </v-col>
        </v-row>
      </v-card>
      <v-row>
        <v-col>
          <v-card title="References" class="small-card-title" flat>
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
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import {
  $ExternalBioMatOccurrenceInput,
  $ExternalBioMatUpdate,
  BioMaterialWithDetails,
  DateWithPrecision,
  ExternalBioMatOccurrenceInput,
  ExternalBioMatUpdate,
  IdentificationInput
} from '@/api'
import {
  createExternalBioMatMutation,
  updateExternalBioMatMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import ArticlesPicker from '../references/ArticlesPicker.vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import ExtBioMatQuantityPicker from './ExtBioMatQuantityPicker.vue'
import IdentificationFormFields from './IdentificationFormFields.vue'
import SitePicker from '../sites/SiteAutocomplete.vue'
import SiteSelectorCard from '../sites/SiteSelectorCard.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<BioMaterialWithDetails>()

const initial: Partial<Omit<ExternalBioMatOccurrenceInput, 'identification'>> & {
  identification: Partial<IdentificationInput>
} = {
  identification: { identified_on: { date: {}, precision: 'Day' } }
}

function updateTransformer({
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
    published_in: published_in?.map(({ code, original }) => ({ code, original })) ?? null
  }
}

const create = {
  mutation: createExternalBioMatMutation,
  schema: $ExternalBioMatOccurrenceInput
}

const update = {
  mutation: updateExternalBioMatMutation,
  schema: $ExternalBioMatUpdate,
  itemID: ({ code }: BioMaterialWithDetails) => ({ code })
}
</script>

<style scoped lang="scss"></style>
