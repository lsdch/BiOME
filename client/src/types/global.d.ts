import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { Replace } from 'ts-toolbelt/out/Object/Replace'
import { Ref } from 'vue'
import { DataTableSortItem } from 'vuetify'
import { VDataTable } from 'vuetify/components'
import { VIcon } from 'vuetify/components/VIcon'
import { DataTableItem } from 'vuetify/lib/components/VDataTable/types.mjs'

export {}

type UnwrapReadonlyArray<A> = A extends Readonly<Array<infer I>> ? I : never
type ReadonlyHeaders = VDataTable['$props']['headers']

declare global {
  type UUID = string & {}

  type IconValue = VIcon['$props']['icon']

  // DataTables

  type DataTableHeader = UnwrapReadonlyArray<ReadonlyHeaders>

  type HeaderDefinitionFor<Item extends {} = Unknown, RowItem extends {} = Unknown> = Omit<
    DataTableHeader,
    'filter'
  > & {
    // Allow filtering using any value type instead of string only
    // See original definition of FilterFunction type:
    // https://github.com/vuetifyjs/vuetify/blob/21241e1762734f639b4ee421e00735d3754181c8/packages/vuetify/src/composables/filter.ts#L19-L19
    readonly filter?: (value: Item, query: string, item: DataTableItem<RowItem>) => boolean
    hide?: Ref<boolean>
  }

  type CRUDTableHeader<Item extends {} = Unknown> = Overwrite<
    HeaderDefinitionFor<Exclude<keyof Item, '$schema'>>,
    {
      key?: Exclude<keyof Item, '$schema'> | DataTableHeader['key']
      readonly filter?: (value: any, query: string, item: DataTableItem<Item>) => boolean
      value?: Exclude<keyof Item, '$schema'> | DataTableHeader['value']
    }
  >
  type CRUDTableHeaders<Item extends object> = {
    [K in Exclude<keyof Item, '$schema'>]: CRUDTableHeader<Item> & { key: K }
  }[Exclude<keyof Item, '$schema'>][]

  type SortItem<K = string> = OverWrite<DataTableSortItem, { key: K }>

  // Type wrangling

  /**
   * Build the union of all paths in an object type
   */
  type ObjectPaths<T extends Record<string, any>> = {
    [K in keyof T]-?: T[K] extends Record<string, any> ? `${K}.${ObjectPaths<T[K]>}` : `${K}`
  }[keyof T]

  type DeepPartial<T extends {}> = {
    [K in keyof T]?: DeepPartial<T[K]>
  }

  type PartialTips<T extends {}> = {
    [K in keyof T]?: T[K] extends {} ? PartialTips<T[K]> : T[K]
  }

  type Multiplable<T, Multiple extends boolean> = true extends Multiple ? T[] : T
}
