<template>
  <FormDialog v-model="dialog" :title="`${mode} gene`" :loading @submit="submit">
    <v-container fluid>
      <v-row>
        <v-col>
          <v-text-field label="Label" v-model="model.label" v-bind="field('label')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field label="Code" v-model="model.code" v-bind="field('code')" />
        </v-col>
        <v-col>
          <v-switch
            label="MOTU delimiter"
            v-model="model.is_MOTU_delimiter"
            color="primary"
            v-bind="field('is_MOTU_delimiter')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            label="Description"
            v-model="model.description"
            v-bind="field('description')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { Gene, GeneInput, GeneUpdate, $GeneInput, $GeneUpdate, SequencesService } from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed, useToggle } from '@vueuse/core'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Gene>>()
const emit = defineEmits<FormEmits<Gene>>()

const initial: GeneInput = {
  code: '',
  label: '',
  is_MOTU_delimiter: false
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ code, label, description, is_MOTU_delimiter }): GeneUpdate {
    return { code, label, description, is_MOTU_delimiter }
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $GeneInput : $GeneUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: SequencesService.createGene,
    edit: ({ code }, body) => SequencesService.updateGene({ path: { code }, body })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
