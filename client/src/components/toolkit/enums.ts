export function enumAsString(v: string) {
  const lowercase = v.split(/(?=[A-Z])/).join(' ').toLowerCase()
  return lowercase[0].toUpperCase() + lowercase.slice(1)
}


// Extract string union type variants as array
// https://stackoverflow.com/a/70694878/12421092
type ValueOf<T> = T[keyof T];

type NonEmptyArray<T> = [T, ...T[]]

type MustInclude<T, U extends T[]> = [T] extends [ValueOf<U>] ? U : never;

export function stringUnionToArray<T>() {
  return <U extends NonEmptyArray<T>>(...elements: MustInclude<T, U>) => elements;
}
