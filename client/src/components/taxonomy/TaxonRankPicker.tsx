import { $TaxonRank, TaxonRank } from '@/api'
import { VSelect } from 'vuetify/components'

export type TaxonRankPickerProps = VSelect['$props'] & {
  modelValue?: TaxonRank
  'onUpdate:ModelValue'?: (value: TaxonRank) => void
}

export function TaxonRankPicker({ modelValue, ...props }: TaxonRankPickerProps) {
  return <v-select items={$TaxonRank.enum} model-value={modelValue} {...props} />
}

export default TaxonRankPicker
