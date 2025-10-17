import { OrgKind } from '@/api'
import { SetupContext } from 'vue'
import { VChip } from 'vuetify/components'

export type OrgKindChipProps = {
  kind: OrgKind
  label?: string
  hideLabel?: boolean
} & VChip['$props']

export function OrgKindChip(
  { kind, label, hideLabel, ...chipProps }: OrgKindChipProps,
  { attrs, slots }: SetupContext
) {
  return (
    <v-chip
      title={label}
      label
      variant="outlined"
      class={[attrs['class'], 'px-1']}
      {...{ ...OrgKind.props[kind], ...chipProps }}
    >
      {{
        prepend: () => <v-icon {...OrgKind.props[kind]} size="small" />,
        default: () => (hideLabel ? null : (slots.default?.() ?? label))
      }}
    </v-chip>
  )
}

export default OrgKindChip
