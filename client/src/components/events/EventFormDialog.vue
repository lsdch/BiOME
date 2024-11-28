<template>
  <FormDialog :title v-bind="$attrs">
    <v-container>
      <v-row>
        <v-col>
          <DateWithPrecisionField v-model="model.performed_on" />
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
import { Event, EventInput, EventUpdate, SiteItem } from '@/api'
import PersonPicker from '../people/PersonPicker.vue'
import { FormProps, useForm } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import DateWithPrecisionField from './DateWithPrecisionField.vue'
import ProgramPicker from './ProgramPicker.vue'
import { computed } from 'vue'

const props = defineProps<FormProps<Event> & { site: SiteItem }>()

const initial: EventInput = {
  performed_by: [],
  performed_on: { precision: 'Day', date: undefined },
  programs: []
}

const title = computed(() => {
  return mode.value === 'Create'
    ? `New event at ${props.site.name}`
    : `Update event at ${props.site.name}`
})

const { model, mode } = useForm(props, {
  initial,
  updateTransformer({ performed_on, performed_by, programs }): EventUpdate {
    return {
      performed_on,
      performed_by: performed_by.map(({ alias }) => alias),
      programs: programs?.map(({ code }) => code)
    }
  }
})
</script>

<style scoped lang="scss"></style>
