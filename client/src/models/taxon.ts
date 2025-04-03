import { Taxon, TaxonInput, TaxonUpdate } from "@/api";
import { Optional } from "ts-toolbelt/out/Object/Optional";
import { Reactive } from "vue";
import { reactive } from "vue";

export type TaxonFormModel = Partial<TaxonInput>

export function initialModel(): Reactive<TaxonFormModel> {
  return reactive({})
}

export function fromTaxon({ id, $schema, meta, ...rest }: Taxon): TaxonFormModel {
  return rest satisfies TaxonFormModel
}

export function toRequestBody({ name, code, rank, parent, status, ...model }: TaxonFormModel): TaxonInput {
  return {
    ...model,
    name: name!,
    code: code!,
    parent: parent!,
    status: status!,
    rank: rank!
  } satisfies TaxonUpdate
}