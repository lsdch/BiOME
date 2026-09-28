import { Country } from '@/api'
import { FunctionalComponent } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { VChip } from 'vuetify/components'

type CountryChipProps = {
  country: Country
} & ComponentProps<typeof VChip>

export const CountryChip: FunctionalComponent<CountryChipProps> = (
  { country, ...props },
  context
) => {
  return (
    <v-tooltip>
      {{
        activator: ({ props: p }: { props: Record<string, any> }) => (
          <v-chip text={country.code} {...{ ...props, ...p }} class="font-monospace" />
        ),
        default: () => country.name
      }}
    </v-tooltip>
  )
}

export default CountryChip
