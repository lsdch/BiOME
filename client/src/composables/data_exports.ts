import { ref } from "vue"

export type Delimiter = '\t' | ',' | ';'

export type QuoteChar = '"' | "'"

export type ExportOptions = {
  delimiter: Delimiter
  quotes: boolean
  quoteChar: QuoteChar
}

export function useExportOptions() {

  const options = ref<ExportOptions>(defaults())

  function defaults(): ExportOptions {
    return {
      delimiter: '\t',
      quotes: true,
      quoteChar: '"'
    }
  }

  function reset() {
    options.value = defaults()
  }

  return { options, defaults, reset }
}
