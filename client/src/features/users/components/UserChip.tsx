import {  User, UserRole } from '@/api'
import { getRoleIcon } from '@/components/icons/UserRoleIcon'
import { VChip } from 'vuetify/components'

export type UserChipProps = {
  user: User
  short?: boolean
} & VChip['$props']

export function UserChip({ user, short, ...chipProps }: UserChipProps) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip
            text={short ? `${user.first_name[0]}. ${user.last_name}` : user.full_name}
            {...{ ...props, ...chipProps }}
          >
            {{
              prepend: user.role
                ? () => (
                    <v-icon
                      icon="mdi-circle-medium"
                      color={getRoleIcon(user.role).color}
                      size="small"
                    />
                  )
                : undefined
            }}
          </v-chip>
        ),
        default: () => (
          <v-card
            title={user.full_name}
            subtitle={user.login}
            class="bg-surface-light small-card-title"
            density="compact"
          >
            {{
              default: () =>
                user.organization ? <v-chip text={user.organization} size="small" /> : undefined,
              prepend: () => (
                <div class="d-flex flex-column align-center mr-2">
                  <UserRole.Icon role={user.role} size="small" />
                </div>
              )
            }}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default UserChip
