import { TaxonRank, TaxonStatus } from "@/api";
import { stringUnionToArray } from "../toolkit/enums";

export const taxonStatusOptions = stringUnionToArray<TaxonStatus>()('Accepted', "Synonym", "Unclassified")

export const taxonRankOptions = stringUnionToArray<TaxonRank>()('Kingdom', 'Phylum', 'Class', "Family", 'Genus', 'Species', 'Subspecies')