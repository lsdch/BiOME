
import { Replace } from 'ts-toolbelt/out/Object/Replace';
import { Ref } from 'vue';
import { DataTableSortItem } from 'vuetify';
import { VDataTable } from 'vuetify/components'
import { VIcon } from 'vuetify/components/VIcon';
import { DataTableItem } from 'vuetify/lib/components/VDataTable/types.mjs';


export { };

type UnwrapReadonlyArray<A> = A extends Readonly<Array<infer I>> ? I : never
type ReadonlyHeaders = VDataTable['$props']['headers']

declare global {
  type IconValue = VIcon['$props']['icon']

  // DataTables

  type DataTableHeader = UnwrapReadonlyArray<ReadonlyHeaders>

  type CRUDTableHeader<Item extends {} = Unknown> = Omit<DataTableHeader, 'filter' | 'key'> & {
    // Allow filtering using any value type instead of string only
    // See original definition of FilterFunction type:
    // https://github.com/vuetifyjs/vuetify/blob/21241e1762734f639b4ee421e00735d3754181c8/packages/vuetify/src/composables/filter.ts#L19-L19
    readonly filter?: (value: any, query: string, item: DataTableItem<Item>) => boolean
    key?: Exclude<(keyof Item), "$schema"> | DataTableHeader['key'];
    value?: Exclude<(keyof Item), "$schema"> | DataTableHeader['value'];
    hide?: Ref<boolean>
  };
  type CRUDTableHeaders = CRUDTableHeader[]

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