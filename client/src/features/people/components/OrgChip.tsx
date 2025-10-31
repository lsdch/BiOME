import { OrganisationInner, OrgKind } from '@/api'
import { VCard, VChip } from 'vuetify/components'

export type OrgChipProps = {
  org: OrganisationInner
  cardProps?: VCard['$props']
} & VChip['$props']

export function OrgChip({ org, ...chipProps }: OrgChipProps) {
  return (
    <v-tooltip transition="scale-transition" content-class="pa-0">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip
            text={org.code}
            prepend-icon={OrgKind.icon(org.kind)}
            color={OrgKind.color(org.kind)}
            {...{ ...props, ...chipProps }}
          ></v-chip>
        ),
        default: () => (
          <v-card
            title={org.name}
            class="bg-transparent small-card-title"
            density="compact"
            flat
            text={org.description}
            {...chipProps.cardProps}
          >
            {{
              subtitle: () => (
                <div class="d-flex align-center ga-2">
                  {org.code}
                  <v-chip
                    label
                    variant="outlined"
                    prepend-icon={OrgKind.icon(org.kind)}
                    color={OrgKind.color(org.kind)}
                    text={OrgKind.humanize(org.kind)}
                    size="small"
                  />
                </div>
              )
            }}
          </v-card>
        )
      }}
    </v-tooltip>
  )
}

export default OrgChip
