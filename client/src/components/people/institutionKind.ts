import { InstitutionKind } from "@/api"
import { stringUnionToArray } from "../toolkit/enums"

export function kindIcon(kind?: InstitutionKind) {
  switch (kind) {
    case "Lab":
      return {
        icon: 'mdi-flask',
        color: 'primary'
      }
    case "FoundingAgency":
      return {
        icon: 'mdi-file-certificate',
        color: 'green'
      }
    case "SequencingPlatform":
      return {
        icon: 'mdi-dna',
        color: 'orange'
      }
    case "Other":
      return {
        icon: 'mdi-home-modern',
        color: 'grey'
      }
    default:
      console.error("Unknown institution kind encountered: ", kind)
      return {}
  }
}

export const institutionKindOptions = stringUnionToArray<InstitutionKind>()('Lab', 'FoundingAgency', 'SequencingPlatform', 'Other')