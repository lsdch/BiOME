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
      <FormDialog
        v-model="dialog"
        :title="`${mode} sequence database`"
        :loading="loading.value"
        @submit="submit"
      >
        <v-container fluid>
          <v-row>
            <v-col>
              <v-text-field label="Label" v-model.trim="model.label" v-bind="field('label')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                label="Code"
                v-model.trim="model.code"
                v-bind="field('code')"
                class="input-font-monospace"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                label="Link template"
                v-model.trim="model.link_template"
                v-bind="field('link_template')"
                class="input-font-monospace"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                label="Description"
                v-model.trim="model.description"
                v-bind="field('description')"
              />
            </v-col>
          </v-row>
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $SeqDBInput, $SeqDBUpdate, SeqDb, SeqDbInput, SeqDbUpdate } from '@/api'
import { createSeqDbMutation, updateSeqDbMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import CreateUpdateForm, {
  FormCreateMutation,
  FormUpdateMutation
} from '../toolkit/forms/CreateUpdateForm.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<SeqDb>()

const initial: SeqDbInput = {
  code: '',
  label: ''
}

function updateTransformer({ code, label, description, link_template }: SeqDb): SeqDbUpdate {
  return { code, label, description, link_template }
}

const create: FormCreateMutation<SeqDb, SeqDbInput> = {
  mutation: createSeqDbMutation,
  schema: $SeqDBInput
}

const update = {
  mutation: updateSeqDbMutation,
  schema: $SeqDBUpdate,
  itemID: ({ code }: SeqDb) => ({ code })
}
</script>

<style scoped lang="scss"></style>
