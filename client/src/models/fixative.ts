import { Fixative, FixativeInput, FixativeUpdateParams } from '@/api'
import { reactive, Reactive } from 'vue'

export type FixativeFormModel = FixativeInput | FixativeUpdateParams

export function initialModel(): Reactive<FixativeInput> {
  return reactive({
    code: '',
    name: ''
  })
}

export function fromFixative({ id, $schema, ...rest }: Fixative): FixativeUpdateParams {
  return rest satisfies FixativeFormModel
}
