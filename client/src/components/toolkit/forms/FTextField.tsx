import { FunctionalComponent } from 'vue'
import { ComponentEmit, ComponentSlots } from 'vue-component-type-helpers'
import { VTextField } from 'vuetify/components'

export type FTextFieldProps = VTextField['$props'] & {
  modelModifiers?: VTextField['$props']['modelModifiers'] & { upper?: boolean }
}

export const FTextField: FunctionalComponent<
  FTextFieldProps,
  ComponentEmit<VTextField>,
  ComponentSlots<VTextField>
> = (props, context) => {
  return (
    <VTextField
      {...{ ...props, ...{ ...context } }}
      onUpdate:modelValue={(v) =>
        context.emit('update:modelValue', props.modelModifiers?.upper ? v.toUpperCase() : v)
      }
    />
  )
}

export default FTextField
