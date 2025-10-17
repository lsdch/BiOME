export function endpointGBIF(suffix: string, params?: { query?: Record<string, string | undefined> }) {
  const url = new URL(`https://api.gbif.org/v1/species/${suffix}`)
  if (params?.query) {
    const query = Object.keys(params.query).reduce(
      (acc, key) =>
        params.query![key] === undefined ? { ...acc } : { ...acc, [key]: params.query![key] },
      {},
    );
    url.search = new URLSearchParams(query).toString()
  }
  return url
}