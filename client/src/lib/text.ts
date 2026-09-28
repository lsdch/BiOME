export function pluralize(
  count: number | undefined,
  singular: string,
  plural: string = singular + 's'
): string {
  return (count ?? 0) <= 1 ? singular : plural
}

export function pluralizeWithCount(
  count: number | undefined,
  singular: string,
  { zero, plural }: { plural?: string; zero?: string } = {}
): string {
  return `${count || (zero ?? 0)} ${pluralize(count, singular, plural)}`
}

export function titleCase(value: string): string {
  return value
    .toLocaleLowerCase()
    .replace(
      /\p{L}[\p{L}\p{M}'’-]*/gu,
      (word) => word.charAt(0).toLocaleUpperCase() + word.slice(1)
    )
}
