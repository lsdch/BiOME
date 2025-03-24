<template>
  <CreateUpdateForm
    v-model="item"
    :initial
    :update-transformer
    :create
    :update
    @success="dialog = false"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog v-model="dialog" :title="title(mode)" @submit="submit" :loading="loading.value">
        <v-container>
          <v-row>
            <v-col cols="12" md="8">
              <v-text-field label="Group label" v-model.trim="model.label" v-bind="field('label')">
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
                    class="font-weight-bold"
                    placeholder="Tag name"
                    color="primary"
                    variant="plain"
                    density="compact"
                    v-bind="field('elements', index, 'label')"
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
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import {
  $HabitatGroupInput,
  $HabitatGroupUpdate,
  HabitatGroup,
  HabitatGroupInput,
  HabitatGroupUpdate,
  HabitatInput,
  HabitatRecord
} from '@/api'
import {
  createHabitatGroupMutation,
  updateHabitatGroupMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { Mode } from '../toolkit/forms/form'
import { defineFormCreate, defineFormUpdate, RequestData } from '@/functions/mutations'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<HabitatGroup>()

// const emit = defineEmits<{
//   created: [group: HabitatGroup]
//   updated: [group: HabitatGroup]
// }>()

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

function updateTransformer({ label, depends, exclusive_elements, elements }: HabitatGroup): State {
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

function title(mode: Mode) {
  return mode == 'Create' ? 'Create habitat group' : `Edit habitats: ${item.value!.label}`
}

const create = defineFormCreate(createHabitatGroupMutation(), {
  initial,
  schema: $HabitatGroupInput,
  requestData({ label, exclusive_elements, elements: elts }) {
    const elements = elts.map<HabitatInput>(({ label, description }) => ({ label, description }))
    return {
      body: {
        label,
        exclusive_elements,
        elements
      }
    }
  }
})

const update = defineFormUpdate(updateHabitatGroupMutation(), {
  schema: $HabitatGroupUpdate,
  itemToModel: updateTransformer,
  requestData: ({ label }, model) => ({
    path: { label },
    body: makeUpdateRequestBody(model)
  })
})

function elementsCountHeadline(model: State) {
  const eltCount = model.elements.length
  return `${eltCount} element${eltCount > 1 ? 's' : ''}`
}

function addElement(model: State) {
  model.elements.unshift({ label: '', description: '', operation: 'create' })
}

function removeElement(model: State, index: number) {
  if (model.elements.length <= index) return
  const element = model.elements[index]
  if (element.operation == 'create') {
    model.elements.splice(index, 1)
  } else {
    element.operation = 'delete'
  }
}

function onElementUpdate(model: State, index: number) {
  if (model.elements.length <= index) return
  const element = model.elements[index]
  if (element.operation == 'create') return
  element.operation =
    element.label != element.initial.label || element.description != element.initial.description
      ? 'update'
      : 'keep'
}

function restoreElement(model: State, index: number) {
  if (model.elements.length <= index) return
  const element = model.elements[index]
  if (element.operation == 'create') return
  element.label = element.initial.label
  element.description = element.initial.description
  element.operation = 'keep'
}

// const [loading, toggleLoading] = useToggle(false)

// function makeRequest() {
//   return mode.value === 'Create'
//     ? HabitatsService.createHabitatGroup({
//         body: getCreateRequestBody(model.value)
//       })
//     : HabitatsService.updateHabitatGroup({
//         path: { code: props.edit!.label },
//         body: getUpdateRequestBody(model.value)
//       })
// }

// async function submit() {
//   toggleLoading(true)
//   await makeRequest()
//     .then(errorHandler)
//     .then((habitatGroup) => {
//       if (mode.value == 'Create') {
//         emit('created', habitatGroup)
//       } else {
//         emit('updated', habitatGroup)
//       }
//       dialog.value = false
//     })
//     .finally(() => toggleLoading(false))
// }

function makeCreateRequestBody({
  label,
  exclusive_elements,
  elements
}: State): RequestData<HabitatGroupInput> {
  return {
    body: {
      label,
      exclusive_elements,
      elements: elements.map<HabitatInput>(({ label, description }) => ({ label, description }))
    }
  }
}

function makeUpdateRequestBody({ label, exclusive_elements, elements }: State): HabitatGroupUpdate {
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
