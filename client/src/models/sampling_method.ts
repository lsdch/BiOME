import { SamplingMethod, SamplingMethodInput, SamplingMethodUpdateParams } from '@/api'
import { reactive, Reactive } from 'vue'

export type SamplingMethodFormModel = SamplingMethodInput | SamplingMethodUpdateParams

export function initialModel(): Reactive<SamplingMethodInput> {
  return reactive({
    code: '',
    name: ''
  })
}

export function fromSamplingMethod({
  code,
  name,
  description
}: SamplingMethod): SamplingMethodUpdateParams {
  return { code, name, description }
}
