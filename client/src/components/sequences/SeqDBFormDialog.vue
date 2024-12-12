<template>
  <FormDialog v-model="dialog" :title="`${mode} sequence database`" :loading @submit="submit">
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

<script setup lang="ts">
import { SeqDb, SeqDbInput, SeqDbUpdate, $SeqDBInput, $SeqDBUpdate, SequencesService } from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed, useToggle } from '@vueuse/core'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<SeqDb>>()
const emit = defineEmits<FormEmits<SeqDb>>()

const initial: SeqDbInput = {
  code: '',
  label: ''
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ code, label, description, link_template }): SeqDbUpdate {
    return { code, label, description, link_template }
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $SeqDBInput : $SeqDBUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: SequencesService.createSeqDb,
    edit: ({ code }, body) => SequencesService.updateSeqDb({ path: { code }, body })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
