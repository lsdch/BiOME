import { UserRole } from '@/api'

export function UserRoleChip({ role }: { role?: UserRole }, context: { attrs?: object }) {
  return (
    <v-chip {...context.attrs}>
      <UserRole.Icon role={role} {...{ class: 'mr-1', size: 'small' }} />
      {role}
    </v-chip>
  )
}

export default UserRoleChip
