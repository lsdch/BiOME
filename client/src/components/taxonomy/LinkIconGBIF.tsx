import IconGBIF from '@/components/icons/IconGBIF'
import { HTMLAttributes, IntrinsicElementAttributes } from 'vue'
import { Events } from 'vue'
import { VBtn } from 'vuetify/components'

export type LinkIconGBIFProps = VBtn['$props'] &
  HTMLAttributes & {
    GBIF_ID: number
    size?: string
    tooltip?: boolean
  }

export function LinkIconGBIF(
  { GBIF_ID, size = 'small', tooltip = true, ...props }: LinkIconGBIFProps,
  { emit }: { emit: VBtn['$emit'] }
) {
  return (
    <v-btn
      {...props}
      size={size}
      icon
      title="GBIF record"
      href={`https://www.gbif.org/species/${GBIF_ID}`}
      target="_blank"
    >
      <IconGBIF size="50%" />
    </v-btn>
  )
}

export default LinkIconGBIF
