import { Person, UserRole } from '@/api'
import { getRoleIcon } from '@/components/icons/UserRoleIcon'
import { VChip } from 'vuetify/components'

export type PersonChipProps = {
  person: Person
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
          >
            {{
              prepend: person.role
                ? () => (
                    <v-icon
                      icon="mdi-circle-medium"
                      color={getRoleIcon(person.role).color}
                      size="small"
                    />
                  )
                : undefined
            }}
          </v-chip>
        ),
        default: () => (
          <v-card
            title={person.full_name}
            subtitle={person.alias}
            class="bg-surface-light small-card-title"
            density="compact"
          >
            {{
              default: () =>
                person.organisation ? (
                  <div class="my-1 d-flex align-center gap-1 flex-wrap">
                    <v-icon class="mx-4" icon="mdi-domain" size="small" />
                    {person.organisation}
                  </div>
                ) : undefined,
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
