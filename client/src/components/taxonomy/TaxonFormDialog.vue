<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        :loading="loading.value"
        v-model="dialog"
        title="Create taxon"
        @submit="submit"
        :fullscreen="xs"
      >
        <v-container>
          <v-row>
            <v-col cols="12" sm="6">
              <TaxonPicker
                label="Parent"
                :ranks="['Order', 'Family', 'Genus', 'Species']"
                :readonly="!!parent"
                item-value="code"
                :modelValue="parent?.code"
                @update:modelValue="
                  (parent: Taxon) => {
                    model.parent = parent.code
                    model.rank = TaxonRank.childRank(parent.rank)!
                  }
                "
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                :modelValue="model.parent !== '' ? model.rank : ''"
                label="New descendant rank"
                variant="plain"
                readonly
                append-icon=""
              />
            </v-col>
          </v-row>
          <!-- {{ unindexedErrors }} -->
          <v-row>
            <v-col cols="12" sm="6">
              <v-text-field v-model.trim="model.name" label="Name" v-bind="field('name')" />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model.trim="model.code"
                label="Code"
                v-bind="field('code')"
                :placeholder="generateCode(model)"
                :persistent-placeholder="(model.name?.length ?? 0) > 0"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" sm="6">
              <v-text-field
                label="Authorship (optional)"
                placeholder="e.g. (Linnaeus, 1758)"
                v-bind="field('authorship')"
                v-model.trim="model.authorship"
              />
            </v-col>
            <v-col cols="12" sm="6">
              <StatusPicker v-model="model.status" v-bind="field('status')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                label="Comments (optional)"
                variant="outlined"
                v-model.trim="model.comment"
              ></v-textarea>
            </v-col>
          </v-row>
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'

import {
  $TaxonInput,
  $TaxonUpdate,
  Taxon,
  TaxonInput,
  TaxonRank,
  TaxonUpdate,
  TaxonWithRelatives
} from '@/api'
import { createTaxonMutation, updateTaxonMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import { useDisplay } from 'vuetify'
import { type FormEmits } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import StatusPicker from './StatusPicker.vue'
import TaxonPicker from './TaxonPicker.vue'

const { xs } = useDisplay()

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Taxon>()

const props = defineProps<{ parent?: Taxon }>()
const emit = defineEmits<FormEmits<TaxonWithRelatives>>()

const initial: TaxonInput = {
  name: '',
  parent: props.parent?.name ?? '',
  rank: 'Subspecies',
  status: 'Unclassified',
  authorship: '',
  code: ''
}

function updateTransformer({ $schema, meta, children_count, ...rest }: Taxon): TaxonUpdate {
  return rest
}

const create = defineFormCreate(createTaxonMutation(), {
  initial,
  schema: $TaxonInput
})

const update = defineFormUpdate(updateTaxonMutation(), {
  schema: $TaxonUpdate,
  itemToModel: updateTransformer,
  requestData: ({ code }) => ({ path: { code } })
})

function generateCode(model: TaxonInput | TaxonUpdate) {
  return model.name?.replace(/\s/g, '_')
}
</script>

<style scoped></style>
