import { QuantityRange, Quantity, UserRole } from '@/api'

function quantityIndicator(quantity: QuantityRange): string {
  if (quantity.lower === quantity.upper) {
    return quantity.lower === 1 ? 'One' : `${quantity.lower}`
  }
  switch (quantity.lower) {
    case 2:
      if (quantity.upper === 5) return 'Few (2-5)'
      break
    case 5:
      if (quantity.upper === 20) return 'Several (5-20)'
      break
    case 20:
      if (quantity.upper === 100) return 'Many (20-100)'
      break
    case 100:
      if (quantity.upper === 1000) return 'Numerous (100-1,000)'
      break
  }
  return `${quantity.lower} - ${quantity.upper}`
}

export function QuantityChip(
  { quantity }: { quantity: QuantityRange },
  context: { attrs?: object }
) {
  return (
    <v-chip class="font-monospace" text={quantityIndicator(quantity)} {...context.attrs}></v-chip>
  )
}

export default QuantityChip
