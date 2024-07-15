<template>
  <FormDialog :loading="loading" v-model="dialog" title="Create taxon">
    <v-form @submit.prevent="submit" class="pb-5">
      <v-row>
        <v-col cols="12" sm="6">
          <TaxonPicker
            label="Parent"
            :ranks="['Order', 'Family', 'Genus', 'Species']"
            @update:modelValue="
              (parent: Taxon) => {
                model.parent = parent.code
                model.rank = childRank(parent.rank)
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
      <v-row>
        <v-spacer></v-spacer>
        <v-btn color="primary" type="submit" variant="text" :loading="loading">Submit</v-btn>
      </v-row>
    </v-form>
  </FormDialog>
</template>

<script setup lang="ts">
import { $TaxonInput, $TaxonStatus, Taxon, TaxonInput, TaxonRank, TaxonomyService } from '@/api'
import { Ref, computed, ref } from 'vue'
import { useForm, type FormEmits, type FormProps } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { watch } from 'vue'
import { childRank } from './rank'
import TaxonPicker from './TaxonPicker.vue'
import StatusPicker from './StatusPicker.vue'

const dialog = defineModel<boolean>()

type Props = FormProps<Taxon> & { parent?: Taxon }

const props = defineProps<Props>()
const emit = defineEmits<FormEmits<Taxon>>()

const initial: TaxonInput = {
  name: '',
  parent: parent?.name ?? '',
  rank: 'Subspecies',
  status: 'Unclassified',
  authorship: '',
  code: ''
}

const { loading, model, field, errorHandler } = useForm(props, $TaxonInput, { initial })

watch(
  () => props.parent,
  (parent) => {
    if (parent !== undefined) {
      model.value.rank = childRank(parent.rank)
      model.value.parent = parent.code
    }
  }
)

const generatedCode = computed(() => {
  return model.value.name.replace(/\s/g, '_')
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
    ({ rank }) => rank == 'Class' || rank == 'Family' || rank == 'Genus' || rank == 'Species'
  )
})

async function submit() {
  if (props.edit) {
    return TaxonomyService.updateTaxon({ path: { code: props.edit.code }, body: model.value })
  } else {
    return TaxonomyService.createTaxon({ body: model.value })
  }
}
</script>

<style scoped></style>
