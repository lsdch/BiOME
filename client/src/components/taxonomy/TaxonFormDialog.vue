<template>
  <FormDialog :loading="loading" v-model="dialog" title="Create taxon">
    <v-form @submit.prevent="submit" class="pb-5">
      <v-row>
        <v-col cols="12" sm="6">
          <v-select
            :items="allowedRanks"
            v-model="model.rank"
            label="Rank"
            variant="outlined"
            prepend-inner-icon="mdi-order-bool-ascending"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <v-autocomplete
            :label="parentLabel"
            :loading="loading"
            cache-items
            required
            :items="candidateParents"
            item-value="code"
            item-title="name"
            v-model="model.parent"
            variant="outlined"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field v-model="model.name" label="Name" />
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field label="Authorship (optional)" placeholder="e.g. (Linnaeus, 1758)" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea label="Comments (optional)" variant="outlined"></v-textarea>
        </v-col>
      </v-row>
      <v-row>
        <v-spacer></v-spacer>
        <v-btn color="primary" type="submit" variant="text" :loading="loading">Submit</v-btn>
      </v-row>
    </v-form>
  </FormDialog>
</template>

<script setup lang="ts">
import { $TaxonInput, Taxon, TaxonInput, TaxonRank, TaxonomyService } from '@/api'
import { Ref, computed, ref } from 'vue'
import { useForm, type FormEmits, type FormProps } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'

const dialog = defineModel<boolean>()
const allowedRanks: TaxonRank[] = ['Subspecies', 'Species', 'Genus', 'Family']

const props = defineProps<FormProps<Taxon>>()
const emit = defineEmits<FormEmits<Taxon>>()

const initial: TaxonInput = {
  name: '',
  parent: '',
  rank: 'Subspecies',
  status: 'Unclassified',
  authorship: '',
  code: ''
}

const parentLabel = computed((): string => {
  switch (model.value.rank) {
    case 'Family':
      return 'Parent class'
    case 'Genus':
      return 'Parent family'
    case 'Species':
      return 'Parent genus'
    case 'Subspecies':
      return 'Parent species'
    default:
      console.warn('Taxon form: encountered rank not in {Species, Subspecies}')
      return 'Parent'
  }
})

const taxa: Ref<Taxon[]> = ref(
  await TaxonomyService.listTaxa()
    .then(({ data: taxa, error }) => {
      if (error !== undefined) {
        // TODO: better error handling
        console.error('Failed to fetch taxa', error)
        return []
      }
      return taxa
    })
    .finally(() => {
      loading.value = false
    })
)
const candidateParents = computed(() => {
  return taxa.value.filter(
    ({ rank }) =>
      (rank == 'Class' && model.value.rank == 'Family') ||
      (rank == 'Family' && model.value.rank == 'Genus') ||
      (rank == 'Genus' && model.value.rank == 'Species') ||
      (rank == 'Species' && model.value.rank == 'Subspecies')
  )
})

async function submit() {
  if (props.edit) {
    return TaxonomyService.updateTaxon({ path: { code: props.edit.code }, body: model.value })
  } else {
    return TaxonomyService.createTaxon({ body: model.value })
  }
}

const { loading, model, field, errorHandler } = useForm(props, $TaxonInput, { initial })
</script>

<style scoped></style>
