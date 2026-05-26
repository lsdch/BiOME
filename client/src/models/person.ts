import { Person, PersonInput, PersonUpdate } from '@/api'
import { reactive, Reactive } from 'vue'

export type PersonFormModel = PersonInput

export function initialModel(): Reactive<PersonFormModel> {
  return reactive({
    first_name: '',
    last_name: ''
  })
}

export function fromPerson({
  comment,
  contact,
  first_name,
  last_name,
  organisation
}: Person): PersonFormModel {
  return {
    first_name,
    last_name,
    contact,
    comment,
    organisation
  }
}

export function toRequestBody(model: PersonFormModel): PersonInput {
  return model satisfies PersonUpdate
}
