import { BaseOccurrenceInput, OccurrenceInput, SamplingInput } from '@/api'
import { reactive, Reactive } from 'vue'

export type OccurrenceModel = DeepPartial<OccurrenceInput> & {
  sampling: SamplingInput
  useSamplingID?: UUID
}

export function initialModel(): Reactive<OccurrenceModel> {
  return reactive({
    sampling: { site: { coordinates: {} }, performed_on: { date: {}, precision: 'Day' } },
    identification: {}
  })
}
