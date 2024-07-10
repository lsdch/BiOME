import { Taxon, TaxonRank } from "@/api";
import { RemovableRef } from "@vueuse/core";
import { InjectionKey, ref, Ref } from "vue";

export type FoldEvent = {
  state: boolean | undefined
  rank: TaxonRank
}

export const MaxRankInjection = Symbol() as InjectionKey<Ref<TaxonRank>>




const selectedTaxon = ref<Taxon>()

export function useTaxonSelection() {

  function select(taxon: Taxon) {
    selectedTaxon.value = taxon
  }

  function drop() {
    selectedTaxon.value = undefined
  }

  return { select, drop, selected: selectedTaxon }
}