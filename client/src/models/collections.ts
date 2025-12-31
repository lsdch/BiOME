import { Collection, CollectionInput, CollectionUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type CollectionFormModel = CollectionInput | CollectionUpdate

export function initialModel(): Reactive<CollectionInput> {
  return reactive({
    label: '',
    code: ''
  })
}

export function fromCollection({ id, meta, $schema, ...rest }: Collection): CollectionUpdate {
  return rest satisfies CollectionFormModel
}