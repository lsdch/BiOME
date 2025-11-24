import { Article, DateWithPrecisionInput, ExternalOccurrenceInput, Quantity } from "@/api";
import { IdentificationModel } from "@/components/forms/occurrence/IdentificationFormFields.vue";
import { reactive, Reactive } from "vue";

export type ExternalOccurrenceModel = Omit<ExternalOccurrenceInput, "published_in" | "identification" | "quantity"> & {
  identification: IdentificationModel
  published_in?: Article[],
  quantity?: Quantity
}

export function initialModel(): Reactive<ExternalOccurrenceModel> {
  return reactive({
    identification: {
      identified_on: { precision: 'Day', date: {} },
    }
  })
}

export function toRequestData({ identification, ...model }: ExternalOccurrenceModel): ExternalOccurrenceInput {
  return {
    ...model,
    published_in: model.published_in?.map(({ code }) => code),
    identification: {
      taxon: identification.taxon!.name,
      identified_by: identification.identified_by!.alias,
      identified_on: identification.identified_on.precision === 'Unknown' ? undefined : identification.identified_on as DateWithPrecisionInput
    },
    quantity: model.quantity!
  }
}