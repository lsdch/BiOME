import { DateWithPrecision, SiteWithOccurrences } from '@/api'
import { DateTime } from 'luxon'
export function lastSamplingDate(site: SiteWithOccurrences): DateWithPrecision | undefined {
  if (!site.samplings) return undefined
  const dates = site.samplings
    .map((s) => s.date)
    .filter((d): d is DateWithPrecision => !!d && !!d.date)
  return dates.length
    ? dates.reduce((latest, current) =>
        DateTime.fromJSDate(current.date) > DateTime.fromJSDate(latest.date) ? current : latest
      )
    : undefined
}
