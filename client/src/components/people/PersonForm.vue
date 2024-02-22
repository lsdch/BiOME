<template>
  <v-form @submit.prevent="submit" v-model="isValid" validate-on="input">
    <v-container fluid>
      <v-row>
        <v-col cols="12" sm="4">
          <v-text-field
            name="first_name"
            label="First name"
            v-model.trim="person.first_name"
            required
            :error-messages="errorMsgs.first_name"
            :rules="inlineRules([required])"
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
            validate-on="blur"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <v-text-field
            name="last_name"
            label="Last name"
            v-model.trim="person.last_name"
            :error-messages="errorMsgs.last_name"
            :rules="inlineRules([required])"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="Contact (optional)"
            v-model="person.contact"
            prepend-inner-icon="mdi-at"
            :rules="inlineRules([email])"
            validate-on="blur"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-select
            label="Institutions (optional)"
            v-model="person.institutions"
            :items="institutions"
            chips
            multiple
            :item-props="({ code, name }) => ({ title: code, subtitle: name })"
            item-value="code"
            :error-messages="errorMsgs.institutions"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-spacer />
        <v-btn
          :loading="loading"
          color="primary"
          variant="plain"
          type="submit"
          text="Submit"
          :disabled="!isValid"
        />
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
import { Emits, Props, inlineRules, useForm } from '../toolkit/form'
import { required, email } from '@vuelidate/validators'

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
  // const data = sanitizeEmptyStrings(person.value)
  if (props.edit) {
    return PeopleService.updatePerson(props.edit.id, data)
  } else {
    return PeopleService.createperson(data)
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
