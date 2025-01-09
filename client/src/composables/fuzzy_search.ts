
import fuzzysort from "fuzzysort";
import { computed, MaybeRef, unref } from "vue";

export type KeysDeclaration<Item extends {}> = Readonly<Array<keyof Item | { key: keyof Item | string & {}, fn: (obj: Item) => string }>>

/**
 * A composable function that provides fuzzy search functionality for filtering items.
 *
 * @template Item - The type of the items to be filtered.
 *
 * @param keys - An array of keys or objects with a key and a function to extract the string value for fuzzy search.
 * @param term - The search term used for fuzzy matching.
 * @param items - The array of items to be filtered.
 * @param  [options] - Optional configuration for the fuzzy search.
 * @param  [options.limit] - The maximum number of results to return.
 * @param  [options.threshold] - The threshold score for considering a match.
 *
 * @returns An object containing:
 * - `highlight`: A function to highlight the matched parts of the search term in the item.
 * - `filteredItems`: A computed ref containing the filtered items based on the search term.
 * - `keysIndex`: An index mapping of keys to their positions in the keys array.
 *
 * @function
 * @name useFuzzyItemsFilter
 *
 * @example
 * ```typescript
 * const { highlight, filteredItems, keysIndex } = useFuzzyItemsFilter(
 *   [{ key: 'name', fn: (item) => item.name }, 'description'],
 *   searchTerm,
 *   items,
 *   { limit: 10, threshold: -1000 }
 * );
 * ```
 */
export function useFuzzyItemsFilter<Item extends {}>(
  keys: KeysDeclaration<Item>,
  term: MaybeRef<string>,
  items: MaybeRef<Item[]>,
  options?: { limit?: number, threshold?: number }
) {

  type ExtractKeys<T extends readonly (keyof Item | { key: keyof Item | string & {} })[]> = T[number] extends infer U
    ? U extends keyof Item
    ? U
    : U extends { key: keyof Item | string & {} }
    ? U['key']
    : never
    : never

  type Keys = ExtractKeys<typeof keys>

  const keysIndex = Object.fromEntries(
    keys.map((v, i) => [typeof v == 'object' ? v.key : v, i])
  ) as Record<Keys, number>

  function highlight(
    item: Fuzzysort.KeysResult<Item>,
    key: Keys,
    options?: {
      highlightOpen?: string
      highlightClose?: string
      baseValue?: string,
      defaultText?: string
    }
  ) {
    return item.obj[key]
      ? item[keysIndex[key]].highlight(options?.highlightOpen, options?.highlightClose) ||
      item.obj[key]
      : (options?.defaultText ?? `Unknown ${String(key)}`)
  }

  const filteredItems = computed(() =>
    fuzzysort.go<Item>(unref(term), unref(items), {
      all: true,
      keys: keys.map((v) => (typeof v === 'object' ? v.fn : String(v))),
      ...options
    })
  )

  return { highlight, filteredItems, keysIndex }
}