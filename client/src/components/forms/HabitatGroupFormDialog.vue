<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} habitat`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container>
      <v-row>
        <v-col cols="12" md="8">
          <v-text-field label="Group label" v-model.trim="model.label" v-bind="schema('label')">
            <template #append>
              <v-btn
                color="primary"
                prepend-icon="mdi-plus"
                text="Add tag"
                @click="addElement(model)"
                variant="text"
              />
            </template>
          </v-text-field>
        </v-col>
        <v-col cols="12" md="4">
          <v-switch
            v-model="model.exclusive_elements"
            label="Elements exclusive"
            color="primary"
            :messages="
              model.exclusive_elements
                ? 'Combining group elements disallowed.'
                : 'Combining group elements allowed.'
            "
          />
        </v-col>
      </v-row>
      <v-row class="align-center">
        <v-col><v-divider /></v-col>
        <v-col cols="0" class="text-center text-secondary">
          {{ elementsCountHeadline(model) }}
        </v-col>
        <v-col>
          <v-divider />
        </v-col>
      </v-row>
      <v-row class="align-stretch">
        <v-col cols="12" sm="6" md="4" v-for="(habitat, index) in model.elements" :key="index">
          <v-card
            :class="[
              'border-s-lg border-opacity-100 rounded-s-0',
              habitat.operation == 'create' ? 'border-success' : 'border-primary'
            ]"
          >
            <template #append>
              <div class="d-flex flex-column">
                <v-btn
                  v-show="model.elements.length > 1 && habitat.operation != 'delete'"
                  color="error"
                  size="x-small"
                  icon="mdi-close"
                  variant="text"
                  @click="removeElement(model, index)"
                />
                <v-btn
                  v-if="['update', 'delete'].includes(habitat.operation)"
                  icon="mdi-restore"
                  size="x-small"
                  variant="text"
                  color="primary"
                  @click="restoreElement(model, index)"
                />
              </div>
            </template>
            <template #title>
              <v-text-field
                v-if="habitat.operation != 'delete'"
                v-model.trim="habitat.label"
                placeholder="Tag name"
                color="primary"
                variant="plain"
                density="compact"
                v-bind="schema('elements', index, 'label')"
                :class="{ 'font-weight-bold': true }"
                :hint="undefined"
                @input="onElementUpdate(model, index)"
              />
              <span v-else class="font-weight-bold text-error text-decoration-line-through">{{
                habitat.label
              }}</span>
            </template>
            <v-card-text class="bg-surface-light py-3">
              <v-textarea
                v-if="habitat.operation != 'delete'"
                v-model.trim="habitat.description"
                variant="plain"
                placeholder="Description (optional)"
                density="compact"
                auto-grow
                :rows="1"
                hide-details
                @input="onElementUpdate(model, index)"
                v-bind="schema('elements', index, 'description')"
              />
              <span v-else class="text-caption font-weight-thin">Marked to be removed</span>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { $HabitatGroupInput, $HabitatGroupUpdate } from '@/api'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { HabitatModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<HabitatModel.HabitatGroupModel>({
  default: HabitatModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: HabitatModel.HabitatGroupModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $HabitatGroupInput : $HabitatGroupUpdate))

function elementsCountHeadline(model: HabitatModel.HabitatGroupModel) {
  const eltCount = model.elements.length
  return `${eltCount} element${eltCount > 1 ? 's' : ''}`
}

function addElement(model: HabitatModel.HabitatGroupModel) {
  model.elements.unshift({ label: '', description: '', operation: 'create' })
}

function removeElement(model: HabitatModel.HabitatGroupModel, index: number) {
  if (model.elements.length <= index) return
  const element = model.elements[index]
  if (element.operation == 'create') {
    model.elements.splice(index, 1)
  } else {
    element.operation = 'delete'
  }
}

function onElementUpdate(model: HabitatModel.HabitatGroupModel, index: number) {
  if (model.elements.length <= index) return
  const element = model.elements[index]
  if (element.operation == 'create') return
  element.operation =
    element.label != element.initial.label || element.description != element.initial.description
      ? 'update'
      : 'keep'
}

function restoreElement(model: HabitatModel.HabitatGroupModel, index: number) {
  if (model.elements.length <= index) return
  const element = model.elements[index]
  if (element.operation == 'create') return
  element.label = element.initial.label
  element.description = element.initial.description
  element.operation = 'keep'
}
</script>

<style scoped lang="scss"></style>
