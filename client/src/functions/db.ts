/**
 * Checks if a value has an ID property,
 * indicating it represents an existing item from the database.
 */
export function hasID<T extends { id: string }, Other extends {}>(value: T | Other | undefined): value is T {
  return !!value && 'id' in value
}