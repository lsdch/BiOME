import { Article, DateWithPrecision, ExternalBioMatInput, Quantity } from "@/api";
import { IdentificationModel } from "@/components/forms/occurrence/IdentificationFormFields.vue";
import { reactive, Reactive } from "vue";

export type ExternalBiomatModel = Omit<ExternalBioMatInput, "published_in" | "identification" | "quantity"> & {
  identification: IdentificationModel
  published_in?: Article[],
  quantity?: Quantity
}

export function initialModel(): Reactive<ExternalBiomatModel> {
  return reactive({
    identification: {
      identified_on: { precision: 'Day', date: {} },
    }
  })
}

export function toRequestData({ identification, ...model }: ExternalBiomatModel): ExternalBioMatInput {
  return {
    ...model,
    published_in: model.published_in?.map(({ code }) => ({ code })),
    identification: {
      taxon: identification.taxon!.name,
      identified_by: identification.identified_by!.alias,
      identified_on: identification.identified_on
    },
    quantity: model.quantity!
  }
}