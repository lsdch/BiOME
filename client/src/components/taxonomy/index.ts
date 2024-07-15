import { Taxon, TaxonRank } from "@/api";
import { useEventBus } from "@vueuse/core";
import { InjectionKey, ref, Ref } from "vue";
import { childRank, parentRank } from "./rank";



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


const foldState = ref<{ [k in TaxonRank]: boolean | undefined }>({
  Kingdom: true,
  Phylum: true,
  Class: true,
  Order: true,
  Family: true,
  Genus: true,
  Species: true,
  Subspecies: true
})
const { emit: emitFold, on: onFold } = useEventBus<TaxonRank>('fold')
const { emit: emitUnfold, on: onUnfold } = useEventBus<TaxonRank>('unfold')

export function useFoldState() {

  function fold(rank: TaxonRank) {
    const child = childRank(rank)
    if (child) fold(child)
    foldState.value[rank] = false
    emitFold(rank)
  }

  function unfold(rank: TaxonRank) {
    const parent = parentRank(rank)
    if (parent) unfold(parent)
    foldState.value[rank] = true
    emitUnfold(rank)
  }

  function isFolded(rank: TaxonRank) {
    return !foldState.value[rank]
  }

  function toggleFold(rank: TaxonRank) {
    return foldState.value[rank] ? fold(rank) : unfold(rank)
  }

  return { fold, unfold, toggleFold, onFold, onUnfold, isFolded }
}