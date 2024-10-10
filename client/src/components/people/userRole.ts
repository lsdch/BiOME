import { UserRole } from "@/api"

export function roleIcon(role?: UserRole) {
  switch (role) {
    case "Admin":
      return {
        icon: 'mdi-star-cog',
        color: 'red'
      }
    case "Maintainer":
      return {
        icon: 'mdi-star-circle',
        color: 'orange'
      }
    case "Contributor":
      return {
        icon: 'mdi-star',
        color: 'primary'
      }
    case "Visitor":
      return {
        icon: 'mdi-circle-medium',
        color: 'green'
      }
    default:
      return {}
  }
}

export const orderedUserRoles: UserRole[] = ['Visitor', 'Contributor', 'Maintainer', 'Admin'] as const

interface UserStatus {
  role: UserRole
}

export function isGranted(user: UserStatus, role: UserRole) {
  return orderedUserRoles.indexOf(user.role) >= orderedUserRoles.indexOf(role)
}