import { Quantity, UserRole } from '@/api'

function quantityIndicator(quantity: Quantity) {
  switch (quantity) {
    case 'Unknown':
      return '?'
    case 'One':
      return '1'
    case 'Several':
      return '1-10'
    case 'Dozen':
      return '~10'
    case 'Tens':
      return '10x'
    case 'Hundred':
      return '100+'
  }
}

export function QuantityChip({ quantity }: { quantity: Quantity }, context: { attrs?: object }) {
  return <v-chip text={quantityIndicator(quantity)} {...context.attrs}></v-chip>
}

export default QuantityChip
