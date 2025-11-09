import { HabitatRecord, Sampling, SamplingInput, SamplingUpdate } from "@/api";
import { DateWithPrecisionModel, toInput as DateWithPrecisionToInput, fromDateWithPrecision } from "@/models/date_with_precision";
import { reactive, Reactive } from "vue";

export type SamplingFormModel = Omit<{
  [K in keyof SamplingInput]: K extends keyof Sampling ? Sampling[K] : never
}, 'habitats' | 'performed_on'> & {
  habitats?: HabitatRecord[],
  performed_on: DateWithPrecisionModel
}

export function initialModel(): Reactive<SamplingFormModel> {
  return reactive({
    performed_on: { precision: 'Day', date: {} }
  })
}

export function fromSampling({
  duration,
  access_points,
  fixatives,
  habitats,
  target_taxa,
  comments,
  methods,
  performed_on,
  performed_by,
  performed_by_groups
}: Sampling): SamplingFormModel {
  return {
    duration,
    access_points,
    comments,
    target_taxa,
    habitats,
    fixatives,
    methods,
    performed_on: fromDateWithPrecision(performed_on),
    performed_by: performed_by ?? [],
    performed_by_groups: performed_by_groups ?? [],
  }
}

export function toRequestBody({
  duration,
  access_points,
  comments,
  target_taxa,
  habitats,
  fixatives,
  methods,
  performed_on,
  performed_by,
  performed_by_groups
}: SamplingFormModel): SamplingInput {
  return {
    duration,
    access_points,
    comments,
    target_taxa: target_taxa?.map(({ name }) => name),
    habitats: habitats?.map(({ label }) => label),
    fixatives: fixatives?.map(({ code }) => code),
    methods: methods?.map(({ code }) => code),
    performed_on: DateWithPrecisionToInput(performed_on),
    performed_by: performed_by?.map(({ alias }) => alias),
    performed_by_groups: performed_by_groups?.map(({ code }) => code),
  } satisfies SamplingUpdate
}