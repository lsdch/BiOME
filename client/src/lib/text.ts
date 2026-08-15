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

export function titleCase(str: string): string {
  return str.replace(/\w\S*/g, (txt) => txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase())
}
