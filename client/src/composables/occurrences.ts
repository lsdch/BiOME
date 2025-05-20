import { SiteWithOccurrences } from "@/api"
import { ScaleBinding } from "vue-leaflet-hexbin"

const bindings: Record<'sites' | 'occurrences', ScaleBinding<SiteWithOccurrences>> = {
  sites: (d) => d.length,
  occurrences: (d) => d.reduce((acc, { data }) => acc + data.occurrences.length, 0),
}

export type ScaleBindingSpec = {
  binding?: keyof typeof bindings
  log?: boolean | 10
}

export function useScaleBinding(spec: ScaleBindingSpec | undefined): ScaleBinding<SiteWithOccurrences> | undefined {
  if (!spec) return undefined
  const { binding, log } = spec
  if (!binding) return undefined
  const scaleBinding = bindings[binding]
  switch (log) {
    case true:
      return (d) => Math.log(scaleBinding(d) + 1)
    case 10:
      return (d) => Math.log10(scaleBinding(d) + 1)
    default:
      return scaleBinding
  }
}