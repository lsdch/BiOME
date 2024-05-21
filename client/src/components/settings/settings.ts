import { InstanceSettings, SettingsService } from "@/api";
import { ref } from "vue";

const settings = ref<InstanceSettings>(await SettingsService.instanceSettings())

export function useInstanceSettings() {

  async function reload() {
    settings.value = await SettingsService.instanceSettings()
    return settings.value
  }

  return { settings: settings.value, reload }
}