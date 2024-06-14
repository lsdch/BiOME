
import { VDataTable } from 'vuetify/components'
import { VIcon } from 'vuetify/components/VIcon';


export { };

type UnwrapReadonlyArray<A> = A extends Readonly<Array<infer I>> ? I : never
type ReadonlyHeaders = VDataTable['$props']['headers']

declare global {
  type IconValue = VIcon['$props']['icon']

  // DataTables

  type DataTableHeader = UnwrapReadonlyArray<ReadonlyHeaders>

  type CRUDTableHeader = Omit<DataTableHeader, 'filter'> & {
    // Allow filtering using any value type instead of string only
    // See original definition of FilterFunction type:
    // https://github.com/vuetifyjs/vuetify/blob/21241e1762734f639b4ee421e00735d3754181c8/packages/vuetify/src/composables/filter.ts#L19-L19
    readonly filter?: (value: any, query: string, item?: any) => boolean
  };
  type CRUDTableHeaders = CRUDTableHeader[]

  type SortItem = VDataTable['$props']['sortBy'] extends Readonly<Array<infer T>> | undefined
    ? T
    : never
}