import { ErrorModel } from "@/api";
import { RequestResult } from "@hey-api/client-fetch";
import { onMounted, ref, Ref } from "vue";

export function useFetchItems<T>(fetchItems: () => RequestResult<T[], ErrorModel, false>) {
  const loading = ref<boolean>(true)

  const items = ref<T[]>([])

  async function fetch(): Promise<T[]> {
    loading.value = true
    return fetchItems().then(({ data, error }) => {
      if (error !== undefined) {
        console.error(`Failed to retrieve items using ${fetchItems.name}`, error)
        return []
      }
      return data
    }).finally(() => loading.value = false)
  }

  onMounted(async () => {
    items.value = await fetch()
  })

  return { loading, items, fetch }
}