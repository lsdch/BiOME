import { DatasetCategory } from '@/api'
import { ComponentProps } from 'vue-component-type-helpers'
import { VIcon } from 'vuetify/components'

type Props = ComponentProps<typeof VIcon> & {
  category: DatasetCategory
}

export function DatasetCategoryIcon({ category, ...props }: Props) {
  return <v-icon icon={DatasetCategory.icon(category)} {...props}></v-icon>
}

export default DatasetCategoryIcon
