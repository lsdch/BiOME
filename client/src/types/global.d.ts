
import { VDataTable } from 'vuetify/components'

export { };

type ExtractDataTableHeader<T> = T extends readonly (infer U)[] ? U : never;

declare global {
  // DataTables
  type ReadonlyHeaders = InstanceType<typeof VDataTable>['headers']
  type CRUDTableHeader = ExtractDataTableHeader<InstanceType<typeof VDataTable>['headers']>;
  type CRUDTableHeaders = CRUDTableHeader[]
}