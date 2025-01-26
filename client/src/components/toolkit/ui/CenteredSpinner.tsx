import { ComponentProps } from 'vue-component-type-helpers'
import { VProgressCircular } from 'vuetify/components'

type CenteredSpinnerProps = ComponentProps<typeof VProgressCircular> & { height?: number | string }

export function CenteredSpinner(props: CenteredSpinnerProps) {
  return (
    <v-sheet class="d-flex align-center justify-center" height={props.height}>
      <v-progress-circular indeterminate {...props} />
    </v-sheet>
  )
}

export default CenteredSpinner
