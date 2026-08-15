import { $OccurrenceTypeStatus, OccurrenceTypeStatus } from '@/api'
import { VSelect } from 'vuetify/components'

export type OccurrenceTypeStatusPickerProps<Multiple extends boolean | undefined> = VSelect['$props'] & {
  multiple: Multiple extends true ? true : false
  modelValue?: Multiple extends true ? OccurrenceTypeStatus[] : OccurrenceTypeStatus
  'onUpdate:ModelValue'?: (value: OccurrenceTypeStatus) => void
}

export function OccurrenceTypeStatusPicker<Multiple extends boolean | undefined>({
  modelValue,
  ...props
}: OccurrenceTypeStatusPickerProps<Multiple>) {
  return <v-select items={$OccurrenceTypeStatus.enum} model-value={modelValue} {...props} style={{ 'text-transform': 'capitalize' }} />
}

export default OccurrenceTypeStatusPicker
