import { DateWithPrecision, DateWithPrecisionInput, HabitatRecord, Sampling, SamplingInput, SamplingUpdate } from "@/api";
import { reactive } from "vue";
import { Reactive } from "vue";

export type SamplingFormModel = Omit<{
  [K in keyof SamplingInput]: K extends keyof Sampling ? Sampling[K] : never
}, 'habitats' | 'performed_on'> & {
  habitats?: HabitatRecord[],
  performed_on?: DateWithPrecisionInput
}

export function initialModel(): Reactive<SamplingFormModel> {
  return reactive({
    target: {
      kind: "Taxa"
    },
    performed_on: { precision: 'Day', date: {} }
  })
}

export function fromSampling({
  duration,
  access_points,
  fixatives,
  habitats,
  target,
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
    target: {
      kind: target.kind,
      taxa: target.taxa
    },
    habitats,
    fixatives,
    methods,
    performed_on: DateWithPrecision.toInput(performed_on),
    performed_by: performed_by ?? [],
    performed_by_groups: performed_by_groups ?? [],
  }
}

export function toRequestBody({
  duration,
  access_points,
  comments,
  target,
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
    target: {
      kind: target.kind,
      taxa: target.taxa?.map(({ name }) => name)
    },
    habitats: habitats?.map(({ label }) => label),
    fixatives: fixatives?.map(({ code }) => code),
    methods: methods?.map(({ code }) => code),
    performed_on: performed_on,
    performed_by: performed_by?.map(({ alias }) => alias),
    performed_by_groups: performed_by_groups?.map(({ code }) => code),
  } satisfies SamplingUpdate
}