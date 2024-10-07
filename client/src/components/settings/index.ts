import { InstanceSettings, SettingsService } from "@/api";
import { handleErrors } from "@/api/responses";
import { computed, ref } from "vue";

const settings = ref<InstanceSettings>()
const ICON_PATH = '/api/v1/assets/app_icon.png' as const
const cacheKey = ref(Math.random())

async function reload() {
  settings.value = await SettingsService.instanceSettings()
    .then(handleErrors((err) => {
      console.error("Failed to fetch instance settings:", err)
    }))
  return settings.value!
}

export function initInstanceSettings() {
  return reload()
}

export function useInstanceSettings() {
  const iconImgProps = computed(() => ({
    src: `${ICON_PATH}?cacheKey=${cacheKey.value}`,
    key: cacheKey.value
  }))

  function reloadIcon() {
    cacheKey.value = Math.random()
  }

  return { settings: settings.value!, reload, ICON_PATH, iconImgProps, reloadIcon }
}
