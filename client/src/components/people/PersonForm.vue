<template>
  <v-form @submit.prevent="submit" v-model="isValid" validate-on="input">
    <v-container fluid>
      <PersonFormFields v-model="person" :error-msgs="errorMsgs" />
      <v-row>
        <v-col>
          <v-text-field
            label="Contact (optional)"
            v-model="person.contact"
            prepend-inner-icon="mdi-at"
            :rules="inlineRules([email])"
            validate-on="blur"
            hint="An e-mail address to contact this person"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-autocomplete
            label="Institutions (optional)"
            v-model="person.institutions"
            :items="institutions"
            item-color="primary"
            chips
            closable-chips
            multiple
            :item-props="({ code, name }) => ({ title: code, subtitle: name })"
            item-value="code"
            :error-messages="errorMsgs.institutions"
            prepend-inner-icon="mdi-domain"
          >
            <template v-slot:chip="{ item, props }">
              <InstitutionKindChip :kind="item.raw.kind" v-bind="props" size="x-small">
                {{ item.raw.code }}
              </InstitutionKindChip>
            </template>
          </v-autocomplete>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea v-model="person.comment" variant="outlined" label="Comments (optional)" />
        </v-col>
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
import { email } from '@vuelidate/validators'
import { onMounted, ref, watchEffect } from 'vue'
import { VForm } from 'vuetify/components'
import { Emits, Props, inlineRules, useForm } from '../toolkit/form'
import PersonFormFields from './PersonFormFields.vue'
import InstitutionKindChip from './InstitutionKindChip.vue'

const isValid = ref(null)

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
  if (props.edit) {
    return PeopleService.updatePerson({ id: props.edit.id, requestBody: data })
  } else {
    return PeopleService.createPerson({ requestBody: data })
  }
}

const emit = defineEmits<Emits<Person>>()

const { submit, loading, errorMsgs } = useForm<PersonInput, Person>(props, emit, request)

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
