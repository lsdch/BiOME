import { $TaxonRank, TaxonRank } from '@/api'
import { VSelect } from 'vuetify/components'

export type TaxonRankPickerProps = VSelect['$props'] & {
  modelValue?: TaxonRank
  exclude?: TaxonRank[]
  'onUpdate:ModelValue'?: (value: TaxonRank) => void
}

export function TaxonRankPicker({ modelValue, exclude, ...props }: TaxonRankPickerProps) {
  const items = exclude ? $TaxonRank.enum.filter((t) => !exclude.includes(t)) : $TaxonRank.enum
  return <v-select items={items} model-value={modelValue} {...props} />
}

export default TaxonRankPicker
