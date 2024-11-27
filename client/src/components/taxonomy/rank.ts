import { $TaxonRank, TaxonRank } from "@/api"


const taxonRankOrder = Object.fromEntries($TaxonRank.enum.map((rank, i) => [rank, i])) as { [k in TaxonRank]: number }


export function isDescendant(rank: TaxonRank, from: TaxonRank) {
  return taxonRankOrder[rank] > taxonRankOrder[from]
}

export function isAscendant(rank: TaxonRank, of: TaxonRank) {
  return taxonRankOrder[rank] < taxonRankOrder[of]
}

export function parentRank(rank: TaxonRank): TaxonRank | undefined {
  return $TaxonRank.enum[$TaxonRank.enum.indexOf(rank) - 1]
}

export function childRank(rank: TaxonRank): TaxonRank | undefined {
  return $TaxonRank.enum[$TaxonRank.enum.indexOf(rank) + 1]
}

export function ranksUpTo(rank: TaxonRank): TaxonRank[] {
  return $TaxonRank.enum.slice($TaxonRank.enum.indexOf(rank))
}

export function ranksDownTo(rank: TaxonRank): TaxonRank[] {
  return $TaxonRank.enum.slice(0, $TaxonRank.enum.indexOf(rank) + 1)
}
