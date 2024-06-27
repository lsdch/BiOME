import { LocationService } from "@/api";
import { handleErrors } from "@/api/responses";
import { ref } from "vue";

export async function useAccessPoints() {
  const accessPoints = ref(await fetch())

  function fetch() {
    return LocationService.getAccessPoints().then(handleErrors(err => console.error("Failed to fetch access points", err)))
  }

  async function refresh() {
    accessPoints.value = await fetch()
  }

  return { refresh, accessPoints }
}