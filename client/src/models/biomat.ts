import { Article, DateWithPrecisionInput, OccurrenceInput } from "@/api";
import { IdentificationModel } from "@/components/forms/occurrence/IdentificationFormFields.vue";
import { reactive, Reactive } from "vue";

export type SpecimenQuantityRangeModel = { lower?: number, upper?: number }
export type SpecimenQuantityModel = number | SpecimenQuantityRangeModel

export type BiomatModel = Omit<OccurrenceInput, "published_in" | "identification" | "quantity"> & {
  identification: IdentificationModel
  published_in?: Article[],
  quantity?: SpecimenQuantityModel
}

export function initialModel(): Reactive<BiomatModel> {
  return reactive({
    identification: {
      identified_on: { precision: 'Day', date: {} },
    }
  })
}

export function toRequestData({ identification, ...model }: BiomatModel): OccurrenceInput {
  return {
    ...model,
    published_in: model.published_in?.map(({ code }) => code),
    identification: {
      taxon: identification.taxon!.name,
      identified_by: identification.identified_by!.alias,
      identified_on: identification.identified_on.precision === 'Unknown' ? undefined : identification.identified_on as DateWithPrecisionInput
    },
    quantity: model.quantity ? (
      typeof model.quantity === 'number' ? [model.quantity, model.quantity] : [model.quantity.lower!, model.quantity.upper!])
      : undefined
  }
}