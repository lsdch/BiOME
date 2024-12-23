import { ErrorModel } from "@/api";
import { RequestResult } from "@hey-api/client-fetch";
import { useToggle } from "@vueuse/core";
import { onMounted, ref, Ref } from "vue";

export function useFetchItem<T>(fetchItem: () => RequestResult<T, ErrorModel, false>, options = { immediate: false }) {
  const [loading, toggleLoading] = useToggle(false)

  const item = ref<T>()
  const error = ref<ErrorModel>()

  async function fetch(): Promise<T | undefined> {
    toggleLoading(true)
    return fetchItem().then(({ data, error }) => {
      if (error !== undefined) {
        console.error(`Failed to retrieve item using ${fetchItem.name}`, error)
        return
      }
      return data
    }).finally(() => toggleLoading(false))
  }

  onMounted(async () => {
    if (options.immediate) item.value = await fetch()
  })

  return { loading, item, fetch, error }
}

export function useFetchItems<T>(fetchItems: () => RequestResult<T[], ErrorModel, false>, options = { immediate: false }) {
  const loading = ref<boolean>(true)

  const items = ref<T[]>([])
  const error = ref<ErrorModel>()

  async function fetch(): Promise<T[]> {
    loading.value = true
    return fetchItems().then(({ data, error: err }) => {
      if (err !== undefined) {
        console.error(`Failed to retrieve items using ${fetchItems.name}`, err)
        error.value = err
        return []
      }
      return data
    }).finally(() => loading.value = false)
  }

  onMounted(async () => {
    items.value = await fetch()
  })

  return { loading, items, fetch, error }
}