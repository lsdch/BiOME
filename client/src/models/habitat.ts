import {
  HabitatGroupWithElements,
  HabitatGroupInput,
  HabitatGroupUpdate,
  HabitatInput,
  Habitat
} from '@/api'
import { Reactive, reactive } from 'vue'

export type HabitatMutation = HabitatInput &
  (
    | {
        initial: Habitat
        operation: 'update' | 'delete' | 'keep'
      }
    | { operation: 'create' }
  )

export type HabitatGroupModel = Omit<HabitatGroupInput, 'elements'> & {
  elements: HabitatMutation[]
}

export function initialModel(): Reactive<HabitatGroupModel> {
  return reactive({
    name: '',
    elements: [{ operation: 'create', name: '', description: '' }]
  })
}

export function fromHabitatGroup({
  name,
  parent_id,

  elements
}: HabitatGroupWithElements): HabitatGroupModel {
  return {
    name,
    elements: elements.map((habitat) => {
      const { id, name, description } = habitat
      return {
        id,
        name,
        description,
        initial: habitat,
        operation: 'keep'
      }
    })
  }
}

export function toCreateRequestBody({ elements, ...model }: HabitatGroupModel): HabitatGroupInput {
  return {
    ...model,
    elements: elements.map(({ name, description }) => ({ name, description }))
  }
}

export function toUpdateRequestBody({
  elements,
  name,
  exclusive_elements
}: HabitatGroupModel): HabitatGroupUpdate {
  return {
    name,
    exclusive_elements,
    ...elements.reduce<
      Required<Pick<HabitatGroupUpdate, 'create_elements' | 'update_elements' | 'delete_elements'>>
    >(
      (acc, { name, description, ...e }) => {
        if (e.operation === 'create') {
          acc.create_elements.push({ name, description })
        } else if (e.operation === 'update') {
          acc.update_elements[e.initial.name] = { name, description }
        } else if (e.operation === 'delete') {
          acc.delete_elements.push(e.initial.name)
        }
        return acc
      },
      {
        create_elements: [],
        update_elements: {},
        delete_elements: []
      }
    )
  }
}
