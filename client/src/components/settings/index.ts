import { InstanceSettings, SettingsService } from "@/api";
import { handleErrors } from "@/api/responses";
import { ref } from "vue";

const settings = ref<InstanceSettings>(
  await SettingsService.instanceSettings()
    .then(handleErrors((err) => {
      console.error("Failed to fetch instance settings:", err)
    }))
)

export function useInstanceSettings() {

  const ICON_PATH = '/api/v1/assets/app_icon.png' as const

  async function reload() {
    settings.value = await SettingsService.instanceSettings()
      .then(handleErrors((err) => {
        console.error("Failed to fetch instance settings:", err)
      }))
    return settings.value
  }

  return { settings: settings.value, reload, ICON_PATH }
}
