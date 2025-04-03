import { DataSource, DataSourceInput, DataSourceUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type DataSourceFormModel = DataSourceInput | DataSourceUpdate

export function initialModel(): Reactive<DataSourceInput> {
  return reactive({
    code: '',
    label: '',
    is_MOTU_delimiter: false,
  })
}

export function fromDataSource({ id, $schema, meta, ...rest }: DataSource): DataSourceUpdate {
  return rest satisfies DataSourceFormModel
}

