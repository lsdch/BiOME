import { Gene, GeneInput, GeneUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type GeneFormModel = GeneInput | GeneUpdate

export function initialModel(): Reactive<GeneInput> {
  return reactive({
    code: '',
    label: '',
    is_MOTU_delimiter: false,
  })
}

export function fromGene({ id, $schema, meta, ...rest }: Gene): GeneUpdate {
  return rest satisfies GeneFormModel
}

