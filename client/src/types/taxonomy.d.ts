export type Taxon = {
  ID: string
  GBIF_ID: number
  name: string
  code: string
  status: string
  anchor: boolean
  authorship?: string
  rank: string
  modified?: string
  created?: string
}

