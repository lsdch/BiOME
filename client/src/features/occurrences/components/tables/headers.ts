import { clamp } from "@vueuse/core"
import { Range } from "ts-toolbelt/out/Number/Range"

export type HeaderExtension<T extends {}> = CRUDTableHeader<T> & {
  position: number
}

export function mergeHeaders<T extends {}, E extends {}>(
  baseHeaders: CRUDTableHeader<T>[],
  extensions?: HeaderExtension<E>[]
): CRUDTableHeader<T & E>[] {
  if (!extensions) return baseHeaders
  const newHeaders: CRUDTableHeader<T & E>[] = [...baseHeaders]
  for (const { position, ...header } of extensions) {
    const clampedPosition = clamp(position, 0, baseHeaders.length)
    newHeaders.splice(clampedPosition, 0, header as CRUDTableHeader<T>)
  }
  return newHeaders
}