import { HabitatGroup, HabitatGroupInput, HabitatGroupUpdate, HabitatInput, HabitatRecord } from "@/api"
import { Reactive, reactive } from "vue"

export type HabitatMutation = HabitatInput & (
  | {
    initial: HabitatRecord
    operation: 'update' | 'delete' | 'keep'
  }
  | { operation: 'create' }
)

export type HabitatGroupModel = Omit<HabitatGroupInput, 'elements'> & { elements: HabitatMutation[] }

export function initialModel(): Reactive<HabitatGroupModel> {
  return reactive({
    label: '',
    elements: [],
  })
}

export function fromHabitatGroup({ label, depends, exclusive_elements, elements }: HabitatGroup): HabitatGroupModel {
  return {
    label,
    exclusive_elements,
    elements: elements.map((habitat) => {
      const { id, label, description } = habitat
      return {
        id,
        label,
        description,
        initial: habitat,
        operation: 'keep'
      }
    })
  }
}

export function toCreateRequestBody({
  elements,
  ...model
}: HabitatGroupModel): HabitatGroupInput {
  return {
    ...model,
    elements: elements.map(({ label, description }) => ({ label, description }))
  }
}

export function toUpdateRequestBody({
  elements,
  label,
  exclusive_elements
}: HabitatGroupModel): HabitatGroupUpdate {
  return {
    label,
    exclusive_elements,
    ...elements.reduce<Required<Pick<HabitatGroupUpdate, 'create_tags' | 'update_tags' | 'delete_tags'>>>(
      (acc, { label, description, ...e }) => {
        if (e.operation === 'create') {
          acc.create_tags.push({ label, description })
        } else if (e.operation === 'update') {
          acc.update_tags[e.initial.label] = { label, description }
        } else if (e.operation === 'delete') {
          acc.delete_tags.push(e.initial.label)
        }
        return acc
      },
      {
        create_tags: [],
        update_tags: {},
        delete_tags: [],
      }
    )
  }
}