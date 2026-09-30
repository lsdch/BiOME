import { formatDistance } from '@/lib/distances'
import type { FunctionalComponent } from 'vue'

export const DistanceDisplay: FunctionalComponent<{ distance: number }> = ({ distance }) => (
  <span class="font-monospace">{formatDistance(distance)}</span>
)