import { UserRole } from "@/api"

export function roleIcon(role?: UserRole) {
  switch (role) {
    case "Admin":
      return {
        icon: 'mdi-star-cog',
        color: 'red'
      }
    case "ProjectMember":
      return {
        icon: 'mdi-star-circle',
        color: 'orange'
      }
    case "Contributor":
      return {
        icon: 'mdi-star',
        color: 'primary'
      }
    case "Guest":
      return {
        icon: 'mdi-circle-medium',
        color: 'green'
      }
    default:
      console.error("Unknown user role encountered: ", role)
      return {}
  }
}