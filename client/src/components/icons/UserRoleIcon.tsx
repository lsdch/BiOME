import { UserRole } from '@/api/adapters'
import { VIcon } from 'vuetify/components'

export const RoleIcon : Record<UserRole, { icon: string; color: string }> = {
  admin: {
    icon: 'mdi-star-cog',
    color: 'red'
  },
  maintainer: {
    icon: 'mdi-star-circle',
    color: 'orange'
  },
  contributor: {
    icon: 'mdi-star',
    color: 'primary'
  },
  visitor: {
    icon: 'mdi-circle-medium',
    color: 'green'
  }
}

export function getRoleIcon(role?: UserRole) {
  return role ? RoleIcon[role] : { icon: 'mdi-account', color: '' }
}


export type UserRoleIconProps = {
  role?: UserRole
} & VIcon['$props']

export function UserRoleIcon({ role, ...props }: UserRoleIconProps, context: { attrs?: object }) {
  const { icon, color } = getRoleIcon(role)
  return (
    <v-tooltip>
      {{
        default: () => `${role || 'No user account'}`,
        activator: ({ props }: { props: any }) => (
          <v-icon icon={icon} color={color} {...{ ...props, ...context.attrs }} />
        )
      }}
    </v-tooltip>
  )
}

export default UserRoleIcon
