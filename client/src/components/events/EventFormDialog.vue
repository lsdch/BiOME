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
      <FormDialog
        :title="
          mode.value === 'Create' ? `New event at ${site.name}` : `Update event at ${site.name}`
        "
        v-bind="$attrs"
        @submit="submit"
        :loading="loading.value"
      >
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
import { DateWithPrecision } from '@/api/adapters'
import { createEventMutation, updateEventMutation } from '@/api/gen/@tanstack/vue-query.gen'
import PersonPicker from '../people/PersonPicker.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import DateWithPrecisionField from './DateWithPrecisionField.vue'
import ProgramPicker from './ProgramPicker.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Event>()
defineProps<{ site: SiteItem }>()

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

const create = {
  mutation: createEventMutation,
  schema: $EventInput
}

const update = {
  mutation: updateEventMutation,
  schema: $EventUpdate,
  itemID: ({ code }: Event) => ({ code })
}
</script>

<style scoped lang="scss"></style>
