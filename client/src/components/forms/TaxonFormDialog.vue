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
            :readonly="!!parent"
            item-value="code"
            return-object
            :modelValue="parent"
            @update:modelValue="
              (parent: Taxon | undefined) => {
                model.parent = parent?.code
                model.rank = parent ? TaxonRank.childRank(parent.rank) : undefined
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
          <v-text-field v-model.trim="model.name" label="Name" v-bind="schema('name')" />
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field
            v-model.trim="model.code"
            label="Code"
            v-bind="schema('code')"
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
import { $TaxonInput, $TaxonUpdate, Taxon, TaxonRank } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { TaxonModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'
import StatusPicker from '../taxonomy/StatusPicker.vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<TaxonModel.TaxonFormModel>({
  default: TaxonModel.initialModel
})

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
</script>

<style scoped lang="scss"></style>
