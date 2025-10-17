import { PersonInner, UserRole } from '@/api'
import { VChip } from 'vuetify/components'

export type PersonChipProps = {
  person: PersonInner
  short?: boolean
} & VChip['$props']

export function PersonChip({ person, short, ...chipProps }: PersonChipProps) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip
            text={short ? `${person.first_name[0]}. ${person.last_name}` : person.full_name}
            {...{ ...props, ...chipProps }}
          />
        ),
        default: () => (
          <v-card
            title={person.full_name}
            subtitle={person.alias}
            class="bg-surface-light small-card-title"
            density="compact"
          >
            {{
              prepend: () => (
                <div class="d-flex flex-column align-center mr-2">
                  <UserRole.Icon role={person.role} size="small" />
                </div>
              )
            }}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default PersonChip
