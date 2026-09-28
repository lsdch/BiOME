import {
  CsvDelimiter,
  CsvQuoteChar,
  ExportFormat,
  ExportSamplingsWithOccurrencesData,
  ListOccurrencesParams
} from '@/api'
import { client } from '@/api/gen/client.gen'
import { MappingFilters } from '@/features/cartography/components/layers-manager/map-layers'

export type ExportOptions = {
  format: ExportFormat
  filename: string
  csvOptions: {
    delimiter: CsvDelimiter
    quoteChar: CsvQuoteChar
  }
}

export async function exportOccurrences(options: ExportOptions, query: ListOccurrencesParams) {
  const queryClient = client
  const requestUrl = queryClient.buildUrl<ExportSamplingsWithOccurrencesData>({
    url: '/occurrences/export',
    query: {
      filename: options.filename,
      format: options.format,
      delimiter: options.csvOptions.delimiter,
      quoteChar: options.csvOptions.quoteChar,
      ...query
    }
  })
  window.location.assign(requestUrl)
}
