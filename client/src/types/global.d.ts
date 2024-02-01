
import { VDataTable } from 'vuetify/components'

export { };

declare global {
  // DataTables
  type ReadonlyHeaders = InstanceType<typeof VDataTable>['headers']

}