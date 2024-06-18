<template>
  <FormDialog v-model="dialog" title="Create habitat group" @submit="submit" :loading="loading">
    <v-form @submit.prevent="submit">
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
          {{
            `${model.elements?.length ?? 0} element${(model.elements?.length ?? 0) > 1 ? 's' : ''}`
          }}
        </v-col>
        <v-col class="pa-0">
          <v-divider />
        </v-col>
      </v-row>
      <v-row class="align-stretch">
        <v-col cols="12" sm="6" md="4" v-for="(habitat, index) in model.elements" :key="index">
          <v-card class="border-primary border-s-lg border-opacity-100 rounded-s-0">
            <template #append>
              <v-btn
                v-show="(model.elements?.length ?? 0) > 1"
                color="error"
                size="x-small"
                icon="mdi-close"
                variant="text"
                @click="model.elements?.splice(index, 1)"
              />
            </template>
            <template #title>
              <v-text-field
                v-model.trim="habitat.label"
                class="font-weight-bold"
                placeholder="Tag name"
                color="primary"
                variant="plain"
                v-bind="field('elements', index, 'label')"
                :hint="undefined"
              />
            </template>
            <v-card-text class="bg-surface-light py-3">
              <v-textarea
                v-model.trim="habitat.description"
                variant="plain"
                placeholder="Description (optional)"
                density="compact"
                auto-grow
                :rows="1"
                hide-details
                v-bind="field('elements', index, 'description')"
              />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-form>
  </FormDialog>
</template>

<script setup lang="ts">
import { $HabitatGroupInput, HabitatGroup, HabitatGroupInput, LocationService } from '@/api'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { FormEmits, FormProps, useForm } from '../toolkit/forms/form'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<HabitatGroup>>()
const emit = defineEmits<FormEmits<HabitatGroup>>()

const initial: HabitatGroupInput = {
  label: '',
  exclusive_elements: true,
  elements: [{ label: '', description: '' }]
}

const { field, errorHandler, model, loading } = useForm(props, $HabitatGroupInput, { initial })

function addElement() {
  model.value.elements?.unshift({ label: '', description: '' }) ||
    (model.value.elements = [{ label: '', description: '' }])
}

async function submit() {
  await LocationService.createHabitatGroup({ body: model.value })
    .then(errorHandler)
    .then((created) => {
      emit('success', created)
      dialog.value = false
    })
}
</script>

<style scoped></style>
