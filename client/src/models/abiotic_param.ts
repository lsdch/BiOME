import { AbioticParameter, AbioticParameterInput, AbioticParameterUpdate } from "@/api";
import { Reactive, reactive } from "vue";

export type AbioticParamModel = AbioticParameterInput | AbioticParameterUpdate

export function initialModel(): Reactive<AbioticParameterInput> {
  return reactive({
    label: '',
    code: '',
    unit: ''
  })
}

export function fromAbioticParam({ id, $schema, meta, ...rest }: AbioticParameter): AbioticParameterUpdate {
  return rest satisfies AbioticParamModel
}