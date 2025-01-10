import { UserRole } from '@/api/adapters'

export function roleIcon(role?: UserRole) {
  switch (role) {
    case 'Admin':
      return {
        icon: 'mdi-star-cog',
        color: 'red'
      }
    case 'Maintainer':
      return {
        icon: 'mdi-star-circle',
        color: 'orange'
      }
    case 'Contributor':
      return {
        icon: 'mdi-star',
        color: 'primary'
      }
    case 'Visitor':
      return {
        icon: 'mdi-circle-medium',
        color: 'green'
      }
    default:
      return {}
  }
}

export function UserRoleIcon({ role }: { role?: UserRole }, context: { attrs?: object }) {
  const { icon, color } = roleIcon(role)
  return <v-icon icon={icon} color={color} title={role} {...context.attrs} />
}

export default UserRoleIcon
