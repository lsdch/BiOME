import { $TaxonRank, TaxonRank } from "@/api"


const taxonRankOrder = Object.fromEntries($TaxonRank.enum.map((rank, i) => [rank, i])) as { [k in TaxonRank]: number }


export function isDescendant(rank: TaxonRank, from: TaxonRank) {
  return taxonRankOrder[rank] > taxonRankOrder[from]
}

export function isAscendant(rank: TaxonRank, of: TaxonRank) {
  return taxonRankOrder[rank] < taxonRankOrder[of]
}