import { Taxon, Taxonomy, TaxonRank } from "@/api";
import { createEventHook, useEventBus, useLocalStorage } from "@vueuse/core";
import { nextTick, reactive, Reactive, ref, ToRefs, toRefs } from "vue";


export const maxRankDisplay = useLocalStorage<TaxonRank>('max-taxon-rank', 'Kingdom')

const selectedTaxon = ref<Taxon>()
const selectHook = createEventHook<Taxon>()

export function useTaxonSelection() {

  function select(taxon: Taxon) {
    selectedTaxon.value = taxon
    history.replaceState(null, '', `#${taxon.name}`)
    selectHook.trigger(taxon)
  }

  function drop() {
    selectedTaxon.value = undefined
    history.replaceState(null, '', ``)
  }

  return { select, drop, selected: selectedTaxon, onSelect: selectHook.on }
}

type TaxonFoldState = ToRefs<Reactive<{ expanded: boolean, parent: string | undefined }>>
const taxonFoldState: Record<string, TaxonFoldState> = {}

function foldTaxon(taxonID: string) {
  return taxonFoldState[taxonID].expanded.value = false
}

async function unfoldTaxon(taxonID: string) {
  const state = taxonFoldState[taxonID]
  if (state.parent.value !== undefined) {
    unfoldTaxon(state.parent.value)
    await nextTick()
  }
  return state.expanded.value = true

}

export async function showTaxon(taxon: Taxonomy) {
  if (taxon.parent) await unfoldTaxon(taxon.parent.id)
  return nextTick(() =>
    scrollToTaxon(taxon.name)
  )
}

export function scrollToTaxon(name: string) {
  document.getElementById(name)!.scrollIntoView({ block: 'center' })
}

export function useTaxonFoldState(taxon: Taxonomy, initial?: boolean) {
  if (!(taxon.id in taxonFoldState))
    taxonFoldState[taxon.id] = toRefs(reactive({ expanded: initial ?? true, parent: taxon.parent?.id }))
  const state = taxonFoldState[taxon.id]

  function fold() {
    return foldTaxon(taxon.id)
  }

  function unfold() {
    return unfoldTaxon(taxon.id)
  }

  function toggleFold() {
    return state.expanded.value ? fold() : unfold()
  }

  async function show() {
    if (taxon.parent) await unfoldTaxon(taxon.parent.id)
    nextTick(() =>
      scrollToTaxon(taxon.name)
    )
  }


  return { expanded: state.expanded, fold, unfold, show, toggleFold, scrollToTaxon }
}

const rankFoldState = ref<{ [k in TaxonRank]: boolean | undefined }>({
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

export function useRankFoldState() {

  function fold(rank: TaxonRank) {
    const child = TaxonRank.childRank(rank)
    if (child) fold(child)
    rankFoldState.value[rank] = false
    emitFold(rank)
  }

  function unfold(rank: TaxonRank) {
    const parent = TaxonRank.parentRank(rank)
    if (parent) unfold(parent)
    rankFoldState.value[rank] = true
    emitUnfold(rank)
  }

  function isFolded(rank: TaxonRank) {
    return !rankFoldState.value[rank]
  }

  function toggleFold(rank: TaxonRank) {
    return rankFoldState.value[rank] ? fold(rank) : unfold(rank)
  }

  return { fold, unfold, toggleFold, onFold, onUnfold, isFolded }
}