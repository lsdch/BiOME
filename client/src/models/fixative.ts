import { Fixative, FixativeInput, FixativeUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type FixativeFormModel = FixativeInput | FixativeUpdate

export function initialModel(): Reactive<FixativeInput> {
  return reactive({
    code: '',
    label: '',
  })
}

export function fromFixative({ id, $schema, meta, ...rest }: Fixative): FixativeUpdate {
  return rest satisfies FixativeFormModel
}

