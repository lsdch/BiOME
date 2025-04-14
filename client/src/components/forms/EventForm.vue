<template>
  <v-container>
    <v-row>
      <v-col>
        <!-- No schema binding, component enforces constraints on its own -->
        <DateWithPrecisionField v-model="model.performed_on" />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <PersonPicker
          label="Performed by"
          v-model="model.performed_by"
          multiple
          return-object
          v-bind="schema('performed_by')"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <OrganisationPicker
          label="Performed by group(s)"
          v-model="model.performed_by_groups"
          multiple
          return-object
          v-bind="schema('performed_by_groups')"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { $EventInput, $EventUpdate } from '@/api'
import PersonPicker from '@/components/people/PersonPicker.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { EventModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'
import OrganisationPicker from '../people/OrganisationPicker.vue'
import DateWithPrecisionField from '../toolkit/forms/DateWithPrecisionField.vue'

const model = defineModel<EventModel.EventModel>({
  default: EventModel.initialModel
})

const { mode = 'Create' } = defineProps<FormProps>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $EventInput : $EventUpdate))
</script>

<style scoped lang="scss"></style>
