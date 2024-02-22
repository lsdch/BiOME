<template>
  <v-form @submit.prevent="submit" class="pb-5">
    <v-row>
      <v-col cols="12" sm="6">
        <v-select :items="allowedRanks" v-model="taxon.rank" label="Rank" variant="outlined" />
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
          v-model="taxon.parent"
          variant="outlined"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6">
        <v-text-field v-model="taxon.name" label="Name" />
      </v-col>
      <v-col cols="12" sm="6">
        <v-text-field label="Authorship (optional)" />
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
</template>

<script setup lang="ts">
import { TaxonInput, TaxonRank, TaxonWithRelatives, TaxonomyService } from '@/api'
import { Ref, computed, onMounted, ref } from 'vue'
import { useForm, Props, Emits } from '../toolkit/form'

const allowedRanks: TaxonRank[] = ['Species', 'Subspecies']

const props = defineProps<Props<TaxonWithRelatives>>()
const emit = defineEmits<Emits<TaxonWithRelatives>>()

const taxon: Ref<TaxonInput> = ref({
  name: '',
  parent: '',
  rank: 'Species',
  status: 'Unclassified',
  authorship: '',
  code: ''
})

const parentLabel = computed((): string => {
  switch (taxon.value.rank) {
    case 'Species':
      return 'Parent genus'
    case 'Subspecies':
      return 'Parent species'
    default:
      console.warn('Taxon form: encountered rank not in {Species, Subspecies}')
      return 'Parent'
  }
})

const taxa: Ref<TaxonWithRelatives[]> = ref([])
const candidateParents = computed(() => {
  return taxa.value.filter(
    ({ rank }) =>
      (rank == 'Genus' && taxon.value.rank == 'Species') ||
      (rank == 'Species' && taxon.value.rank == 'Subspecies')
  )
})
onMounted(async () => {
  taxa.value = await TaxonomyService.taxonomyList()
  loading.value = false
})

async function request() {
  if (props.edit) {
    return TaxonomyService.updateTaxon(props.edit.code, taxon.value)
  } else {
    return TaxonomyService.createTaxon(taxon.value)
  }
}

const { loading, submit } = useForm(props, emit, request)
</script>

<style scoped></style>
