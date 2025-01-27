import { SetupContext, SlotsType, VNodeChild } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { VProgressCircular } from 'vuetify/components'

type CenteredSpinnerProps = ComponentProps<typeof VProgressCircular> & {
  height?: number | string
  text?: string
}

export function CenteredSpinner(
  props: CenteredSpinnerProps,
  context: SetupContext<unknown, SlotsType<{ prepend: () => VNodeChild }>>
) {
  return (
    <v-sheet class="d-flex align-center justify-center" height={props.height}>
      {context.slots.prepend?.()}
      <v-progress-circular indeterminate {...props} class="mx-1" />
      {props.text}
    </v-sheet>
  )
}

export default CenteredSpinner
