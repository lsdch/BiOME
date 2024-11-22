<template>
  <FormDialog v-model="dialog" :title @submit="submit" :loading="loading">
    <v-container>
      <v-row>
        <v-col cols="12" md="8">
          <v-text-field label="Group label" v-model.trim="model.label" v-bind="field('label')">
            <template #append>
              <v-btn
                color="primary"
                prepend-icon="mdi-plus"
                text="Add tag"
                @click="addElement"
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
          {{ elementsCountHeadline }}
        </v-col>
        <v-col class="pa-0">
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
                  @click="removeElement(index)"
                />
                <v-btn
                  v-if="['update', 'delete'].includes(habitat.operation)"
                  icon="mdi-restore"
                  size="x-small"
                  variant="text"
                  color="primary"
                  @click="restoreElement(index)"
                />
              </div>
            </template>
            <template #title>
              <v-text-field
                v-if="habitat.operation != 'delete'"
                v-model.trim="habitat.label"
                class="font-weight-bold"
                placeholder="Tag name"
                color="primary"
                variant="plain"
                density="compact"
                v-bind="field('elements', index, 'label')"
                :hint="undefined"
                @input="onElementUpdate(index)"
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
                @input="onElementUpdate(index)"
                v-bind="field('elements', index, 'description')"
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
import {
  $HabitatGroupInput,
  HabitatGroup,
  HabitatGroupInput,
  HabitatGroupUpdate,
  HabitatInput,
  HabitatRecord,
  HabitatsService
} from '@/api'
import { useToggle } from '@vueuse/core'
import { computed } from 'vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { FormProps, useForm, useSchema } from '../toolkit/forms/form'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<HabitatGroup>>()

const emit = defineEmits<{
  created: [group: HabitatGroup]
  updated: [group: HabitatGroup]
}>()

type Element = HabitatInput &
  (
    | {
        initial: HabitatRecord
        operation: 'update' | 'delete' | 'keep'
      }
    | {
        operation: 'create'
      }
  )

type State = Omit<HabitatGroupInput, 'elements'> & { elements: Element[] }

const initial: State = {
  label: '',
  exclusive_elements: true,
  elements: [{ label: '', description: '', operation: 'create' }]
}

const { model, mode } = useForm(props, {
  initial,
  updateTransformer({ label, depends, exclusive_elements, elements }): State {
    return {
      label,
      exclusive_elements,
      elements: elements.map((habitat) => {
        const { id, label, description } = habitat
        return {
          id,
          label,
          description,
          initial: habitat,
          operation: 'keep'
        }
      })
    }
  }
})
const { field, errorHandler } = useSchema($HabitatGroupInput)

const title = computed(() => {
  return mode.value == 'Create' ? 'Create habitat group' : `Edit habitats: ${props.edit?.label}`
})

const elementsCountHeadline = computed(() => {
  const eltCount = model.value.elements.length
  return `${eltCount} element${eltCount > 1 ? 's' : ''}`
})

function addElement() {
  model.value.elements.unshift({ label: '', description: '', operation: 'create' })
}

function removeElement(index: number) {
  if (model.value.elements.length <= index) return
  const element = model.value.elements[index]
  if (element.operation == 'create') {
    model.value.elements.splice(index, 1)
  } else {
    element.operation = 'delete'
  }
}

function onElementUpdate(index: number) {
  if (model.value.elements.length <= index) return
  const element = model.value.elements[index]
  if (element.operation == 'create') return
  element.operation =
    element.label != element.initial.label || element.description != element.initial.description
      ? 'update'
      : 'keep'
}

function restoreElement(index: number) {
  if (model.value.elements.length <= index) return
  const element = model.value.elements[index]
  if (element.operation == 'create') return
  element.label = element.initial.label
  element.description = element.initial.description
  element.operation = 'keep'
}

const [loading, toggleLoading] = useToggle(false)

function makeRequest() {
  return mode.value === 'Create'
    ? HabitatsService.createHabitatGroup({
        body: getCreateRequestBody(model.value)
      })
    : HabitatsService.updateHabitatGroup({
        path: { code: props.edit!.label },
        body: getUpdateRequestBody(model.value)
      })
}

async function submit() {
  toggleLoading(true)
  await makeRequest()
    .then(errorHandler)
    .then((habitatGroup) => {
      if (mode.value == 'Create') {
        emit('created', habitatGroup)
      } else {
        emit('updated', habitatGroup)
      }
      dialog.value = false
    })
    .finally(() => toggleLoading(false))
}

function getCreateRequestBody({ label, exclusive_elements, elements }: State): HabitatGroupInput {
  return {
    label,
    exclusive_elements,
    elements: elements.map(({ label, description }) => ({ label, description }))
  }
}

function getUpdateRequestBody({ label, exclusive_elements, elements }: State) {
  const body: HabitatGroupUpdate = {
    label,
    exclusive_elements,
    delete_tags: [],
    create_tags: [],
    update_tags: {}
  }
  elements.forEach((e) => {
    switch (e.operation) {
      case 'delete':
        body.delete_tags!.push(e.initial.label)
        break
      case 'create':
        body.create_tags!.push({ label: e.label, description: e.description })
        break
      case 'update':
        body.update_tags![e.initial.label] = { label: e.label, description: e.description }
        break
      default:
        break
    }
  })
  return body
}
</script>

<style scoped></style>
