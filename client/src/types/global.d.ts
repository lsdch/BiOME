
import { VDataTable } from 'vuetify/components'

export { };

declare global {
  // DataTables
  type ReadonlyHeaders = InstanceType<typeof VDataTable>['headers']

  type DataTableHeader = { [P in keyof ReadonlyHeaders[number]]: ReadonlyHeaders[number][P] };

}