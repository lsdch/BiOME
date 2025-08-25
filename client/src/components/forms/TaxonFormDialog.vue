<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} taxon`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container>
      <v-row>
        <v-col cols="12" sm="6">
          <TaxonPicker
            label="Parent"
            :ranks="['Order', 'Family', 'Genus', 'Species']"
            item-value="name"
            return-object
            :multiple="false"
            v-model="parent"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <v-select
            label="New taxon rank"
            v-model="model.rank"
            :items="
              parent
                ? [
                    TaxonRank.childRank(parent.rank),
                    ...(parent.rank === 'Genus' ? ['Subgenus'] : [])
                  ]
                : []
            "
            :disabled="parent?.rank !== 'Genus'"
          />
        </v-col>
      </v-row>
      <!-- {{ unindexedErrors }} -->
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field v-model.trim="model.name" label="Name" v-bind="schema('name')" />
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field
            v-model.trim="model.name"
            label="Code"
            v-bind="schema('name')"
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
            v-bind="schema('authorship')"
            v-model.trim="model.authorship"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <StatusPicker v-model="model.status" v-bind="schema('status')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea label="Comments (optional)" variant="outlined" v-model.trim="model.comment" />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { $TaxonInput, $TaxonUpdate, Taxon, TaxonRank } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { TaxonModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'
import StatusPicker from '../taxonomy/StatusPicker.vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import { watch } from 'vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<TaxonModel.TaxonFormModel>({
  default: TaxonModel.initialModel
})

const parent = defineModel<Taxon>('parent')

const { mode = 'Create', ...props } = defineProps<
  { parent?: Taxon } & FormProps & FormDialogProps
>()

const emit = defineEmits<{
  submit: [model: TaxonModel.TaxonFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $TaxonInput : $TaxonUpdate))

function generateCode(model: TaxonModel.TaxonFormModel) {
  return model.name?.replace(/\s/g, '_')
}

watch(
  parent,
  (newParent) => {
    if (newParent) {
      model.value.parent = newParent.name
      model.value.rank = TaxonRank.childRank(newParent.rank)
    } else {
      model.value.parent = undefined
      model.value.rank = undefined
    }
  },
  { immediate: true }
)
</script>

<style scoped lang="scss"></style>
