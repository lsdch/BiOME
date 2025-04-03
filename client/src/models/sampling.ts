import { HabitatRecord, Sampling, SamplingInput, SamplingUpdate } from "@/api";
import { reactive } from "vue";
import { Reactive } from "vue";

export type SamplingFormModel = Omit<{
  [K in keyof SamplingInput]: K extends keyof Sampling ? Sampling[K] : never
}, 'habitats'> & {
  habitats?: HabitatRecord[]
}

export function initialModel(): Reactive<SamplingFormModel> {
  return reactive({
    target: {
      kind: "Taxa"
    }
  })
}

export function fromSampling({
  duration,
  access_points,
  fixatives,
  habitats,
  target,
  comments,
  methods
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
  }
}

export function toRequestBody({
  duration,
  access_points,
  comments,
  target,
  habitats,
  fixatives,
  methods
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
    methods: methods?.map(({ code }) => code)
  } satisfies SamplingUpdate
}