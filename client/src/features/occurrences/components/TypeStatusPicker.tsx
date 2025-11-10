import { $TypeStatus, TypeStatus } from '@/api'
import { VSelect } from 'vuetify/components'

export type TypeStatusPickerProps<Multiple extends boolean | undefined> = VSelect['$props'] & {
  multiple: Multiple extends true ? true : false
  modelValue?: Multiple extends true ? TypeStatus[] : TypeStatus
  'onUpdate:ModelValue'?: (value: TypeStatus) => void
}

export function TypeStatusPicker<Multiple extends boolean | undefined>({
  modelValue,
  ...props
}: TypeStatusPickerProps<Multiple>) {
  return <v-select items={$TypeStatus.enum} model-value={modelValue} {...props} />
}

export default TypeStatusPicker
