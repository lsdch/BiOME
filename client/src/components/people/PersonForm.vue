<template>
  <v-form @submit.prevent="submit">
    <v-container fluid>
      <v-row class="mb-3">
        <v-col cols="12" sm="4">
          <v-text-field
            name="first_name"
            label="First name"
            v-model.trim="person.first_name"
            required
            :error-messages="errorMsgs.first_name"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <v-text-field
            name="middle_names"
            label="Middle name(s)"
            hint="Optional"
            clearable
            v-model.trim="person.middle_names"
            :error-messages="errorMsgs.middle_names"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <v-text-field
            name="last_name"
            label="Last name"
            v-model.trim="person.last_name"
            :error-messages="errorMsgs.last_name"
          />
        </v-col>
      </v-row>
      <v-row class="mb-3">
        <v-select
          label="Institutions"
          v-model="person.institutions"
          :items="institutions"
          chips
          multiple
          :item-props="({ code, name }) => ({ title: code, subtitle: name })"
          item-value="code"
          :error-messages="errorMsgs.institutions"
        />
      </v-row>
      <v-row>
        <v-spacer />
        <v-btn :loading="loading" color="primary" variant="plain" type="submit" text="Submit" />
      </v-row>
    </v-container>
  </v-form>
</template>

<script lang="ts">
const DEFAULT: PersonInput = {
  first_name: '',
  last_name: '',
  institutions: []
}
</script>

<script setup lang="ts">
import { Institution, PeopleService, Person, PersonInput } from '@/api'
import { onMounted, ref, watchEffect } from 'vue'
import { VForm } from 'vuetify/components'
import { Emits, Props, useForm } from '../toolkit/form'

const props = defineProps<Props<Person>>()

const person = ref(DEFAULT)

watchEffect(() => {
  if (props.edit)
    person.value = {
      ...props.edit,
      institutions: props.edit.institutions?.map(({ code }) => code) ?? []
    }
  else person.value = { ...DEFAULT }
})

function request() {
  const data = person.value
  // const data = sanitizeEmptyStrings(person.value)
  if (props.edit) {
    return PeopleService.updatePerson(props.edit.id, data)
  } else {
    return PeopleService.createperson(data)
  }
}

const emit = defineEmits<Emits<Person>>()

const { submit, loading, errorMsgs } = useForm<PersonInput, Person>(
  props,
  emit,
  request
)

/**
 * List of known institutions from the DB
 */
const institutions = ref<Institution[]>([])
onMounted(async () => {
  institutions.value = await PeopleService.listInstitutions()
})

defineExpose({ submit })
</script>

<style scoped></style>
