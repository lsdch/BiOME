import { FunctionalComponent } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { VIcon } from 'vuetify/components'

type SiteStatusIconProps = {
  exists?: boolean
} & ComponentProps<typeof VIcon>

export const SiteStatusIcon: FunctionalComponent<SiteStatusIconProps> = (
  { exists, ...props },
  context
) => {
  return (
    <v-icon
      {...{
        ...(exists
          ? {
              icon: 'mdi-link-variant',
              color: 'primary',
              title: 'Existing site'
            }
          : {
              icon: 'mdi-plus-thick',
              color: 'warning',
              title: 'New site'
            }),
        ...props
      }}
    />
  )
}

export default SiteStatusIcon
