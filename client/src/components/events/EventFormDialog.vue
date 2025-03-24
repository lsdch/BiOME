<template>
  <CreateUpdateForm
    v-model="item"
    :create
    :update
    :local
    @success="dialog = false"
    @save="dialog = false"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="mode === 'Create' ? `New event` : `Update event`"
        v-bind="$attrs"
        @submit="submit"
        :loading="loading.value"
      >
        <template #subtitle>
          <v-chip label :text="site.code" class="font-monospace" prepend-icon="mdi-map-marker" />
        </template>
        <!-- Expose activator slot -->
        <template #activator="slotData">
          <slot name="activator" v-bind="slotData"></slot>
        </template>

        <v-container>
          <v-row>
            <v-col>
              <DateWithPrecisionField v-model="model.performed_on" v-bind="field('performed_on')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <PersonPicker
                label="Performed by"
                v-model="model.performed_by"
                item-value="alias"
                multiple
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ProgramPicker
                v-model="model.programs"
                item-value="code"
                multiple
                chips
                closable-chips
                clearable
              />
            </v-col>
          </v-row>
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $EventInput, $EventUpdate, Event, EventInput, EventUpdate, SiteItem } from '@/api'
import { DateWithPrecision, SiteInput } from '@/api/adapters'
import { createEventMutation, updateEventMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import PersonPicker from '../people/PersonPicker.vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FormDialog, { FormDialogProps } from '../toolkit/forms/FormDialog.vue'
import DateWithPrecisionField from './DateWithPrecisionField.vue'
import ProgramPicker from './ProgramPicker.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Event>()
const { site } = defineProps<
  {
    site: SiteItem | SiteInput
    local?: true
  } & Omit<FormDialogProps, 'loading' | 'fullscreen'>
>()

const initial: EventInput = {
  performed_by: [],
  performed_on: { precision: 'Day', date: {} },
  programs: []
}

function updateTransformer({ performed_on, performed_by, programs }: Event): EventUpdate {
  return {
    performed_on: DateWithPrecision.toInput(performed_on),
    performed_by: performed_by.map(({ alias }) => alias),
    programs: programs?.map(({ code }) => code)
  }
}

const create = defineFormCreate(createEventMutation(), {
  initial,
  schema: $EventInput,
  requestData(model) {
    return {
      body: model,
      path: { code: site.code }
    }
  }
})

const update = defineFormUpdate(updateEventMutation(), {
  schema: $EventUpdate,
  itemToModel: updateTransformer,
  requestData(item) {
    return {
      path: { id: item.id }
    }
  }
})
</script>

<style scoped lang="scss"></style>
