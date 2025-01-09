<template>
  <FormDialog :title v-bind="$attrs" @submit="submit" :loading>
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

<script setup lang="ts">
import {
  $EventInput,
  $EventUpdate,
  Event,
  EventInput,
  EventsService,
  EventUpdate,
  LocationService,
  SiteItem
} from '@/api'
import { DateWithPrecision } from '@/api/adapters'
import { reactiveComputed, useToggle } from '@vueuse/core'
import { computed, toRefs } from 'vue'
import PersonPicker from '../people/PersonPicker.vue'
import { FormEmits, FormProps, useForm, useSchema } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import DateWithPrecisionField from './DateWithPrecisionField.vue'
import ProgramPicker from './ProgramPicker.vue'

const props = defineProps<FormProps<Event> & { site: SiteItem }>()
const emit = defineEmits<FormEmits<Event>>()
const initial: EventInput = {
  performed_by: [],
  performed_on: { precision: 'Day', date: {} },
  programs: []
}

const [loading, toggleLoading] = useToggle(false)

const title = computed(() => {
  return mode.value === 'Create'
    ? `New event at ${props.site.name}`
    : `Update event at ${props.site.name}`
})

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ performed_on, performed_by, programs }): EventUpdate {
    return {
      performed_on: DateWithPrecision.toInput(performed_on),
      performed_by: performed_by.map(({ alias }) => alias),
      programs: programs?.map(({ code }) => code)
    }
  }
})

const { field, errorHandler } = toRefs(
  reactiveComputed(() => useSchema(mode.value === 'Create' ? $EventInput : $EventUpdate))
)

async function submit() {
  toggleLoading(true)
  await makeRequest({
    create: ({ body }) => LocationService.createEvent({ path: { code: props.site.code }, body }),
    edit: ({ id }, body) => EventsService.updateEvent({ path: { id }, body })
  })
    .then(errorHandler.value)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
