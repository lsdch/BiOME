export function formatDistance(meters: number): string {
  if (meters < 1000) {
    return `${meters.toFixed(0)} m`
  } else {
    const km = (meters / 1000).toFixed(2)
    return `${km} km`
  }
}