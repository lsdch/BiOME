<template>
  <FormDialog :loading v-model="dialog" title="Create taxon" @submit="submit" :fullscreen="xs">
    <v-container>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            v-if="parent"
            :modelValue="parent.name"
            :text="parent.name"
            label="Parent"
            readonly
            variant="plain"
          />
          <TaxonPicker
            v-else
            label="Parent"
            :ranks="['Order', 'Family', 'Genus', 'Species']"
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
      {{ unindexedErrors }}
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field v-model.trim="model.name" label="Name" v-bind="field('name')" />
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field
            v-model.trim="model.code"
            label="Code"
            v-bind="field('code')"
            :placeholder="generatedCode"
            :persistent-placeholder="model.name.length > 0"
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

<script setup lang="ts">
import {
  $TaxonInput,
  Taxon,
  TaxonInput,
  TaxonRank,
  TaxonWithRelatives,
  TaxonomyService
} from '@/api'
import { useToggle } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useDisplay } from 'vuetify'
import { useSchema, type FormEmits, type FormProps } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import StatusPicker from './StatusPicker.vue'
import TaxonPicker from './TaxonPicker.vue'

const { xs } = useDisplay()

const dialog = defineModel<boolean>()

type Props = FormProps<Taxon> & { parent?: Taxon }

const props = defineProps<Props>()
const emit = defineEmits<FormEmits<TaxonWithRelatives>>()

const initial: TaxonInput = {
  name: '',
  parent: props.parent?.name ?? '',
  rank: 'Subspecies',
  status: 'Unclassified',
  authorship: '',
  code: ''
}

const model = ref(initial)

const { field, errorHandler, unindexedErrors } = useSchema($TaxonInput)

watch(
  () => props.parent,
  (parent) => {
    if (parent !== undefined) {
      model.value.rank = TaxonRank.childRank(parent.rank)!
      model.value.parent = parent.code
    }
  },
  { immediate: true }
)

const generatedCode = computed(() => {
  return model.value.name.replace(/\s/g, '_')
})

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  const data = await TaxonomyService.createTaxon({ body: model.value })
    // .then(handleErrors(console.log))
    .then(errorHandler)
    .finally(() => toggleLoading(false))
  emit('success', data)
  dialog.value = false
}
</script>

<style scoped></style>
