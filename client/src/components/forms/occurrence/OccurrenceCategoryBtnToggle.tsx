import { OccurrenceCategory } from '@/api'
import { ComponentProps } from 'vue-component-type-helpers'
import { VBtn } from 'vuetify/components'

type CategoryBtnToggleProps = ComponentProps<typeof VBtn> & {
  modelValue?: OccurrenceCategory
}

export function OccurrenceCategoryBtnToggle(
  { modelValue, ...props }: CategoryBtnToggleProps,
  { emit }: { emit: (event: 'update:modelValue', value: OccurrenceCategory) => void }
) {
  return (
    <v-input
      model-value={modelValue}
      onUpdate:modelValue={(v: OccurrenceCategory) => emit('update:modelValue', v)}
      rules={[(v: OccurrenceCategory | undefined) => !!v || 'Category is required']}
      hide-details
    >
      <v-btn-toggle
        model-value={modelValue}
        onUpdate:modelValue={(v: OccurrenceCategory) => emit('update:modelValue', v)}
        mandatory
        divided
        rounded="md"
        variant="outlined"
      >
        <v-btn
          value="Internal"
          text="Internal"
          prepend-icon={OccurrenceCategory.icon('Internal')}
          color={OccurrenceCategory.props.Internal.color}
          {...props}
        />
        <v-btn
          value="External"
          text="External"
          prepend-icon={OccurrenceCategory.icon('External')}
          color={OccurrenceCategory.props.External.color}
          {...props}
        />
      </v-btn-toggle>
    </v-input>
  )
}
