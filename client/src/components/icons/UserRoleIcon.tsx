import { UserRole } from '@/api/adapters'

interface Props {
  role?: UserRole
}

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

export const UserRoleIcon = (props: Props, attrs?: object) => (
  <v-icon
    {...{
      title: props.role,
      ...roleIcon(props.role),
      ...attrs
    }}
  />
)

export default UserRoleIcon
