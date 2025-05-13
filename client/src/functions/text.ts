export function pluralize(count: number, singular: string, plural: string = singular + "s"): string {
  return count === 1 ? singular : plural;
}

export function pluralizeWithCount(count: number, singular: string, plural: string = singular + "s"): string {
  return `${count} ${pluralize(count, singular, plural)}`;
}