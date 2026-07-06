import { getInstanceSettingsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { injectSettings } from '@/lib/injection'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const ICON_PATH = '/api/v1/assets/app_icon.png' as const
const cacheKey = ref(Math.random())

export function useInstanceSettings() {
  const { data, error, isPending, isSuccess, refetch, suspense } = useQuery({
    ...getInstanceSettingsOptions(),
    initialData: () => {
      return injectSettings()
    }
  })

  const iconImgProps = computed(() => ({
    src: `${ICON_PATH}?cacheKey=${cacheKey.value}`,
    key: cacheKey.value
  }))

  function reloadIcon() {
    cacheKey.value = Math.random()
  }
  return { instance: data, reload: refetch, ICON_PATH, iconImgProps, reloadIcon, isPending, error }
}
