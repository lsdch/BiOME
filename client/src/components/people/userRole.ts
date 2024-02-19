import { UserRole } from "@/api"

export function roleIcon(role?: UserRole) {
  switch (role) {
    case UserRole.Admin:
      return {
        icon: 'mdi-star-cog',
        color: 'red'
      }
    case UserRole.ProjectMember:
      return {
        icon: 'mdi-star-circle',
        color: 'orange'
      }
    case UserRole.Contributor:
      return {
        icon: 'mdi-star',
        color: 'primary'
      }
    case UserRole.Guest:
      return {
        icon: 'mdi-circle-medium',
        color: 'green'
      }
    default:
      return {}
  }
}