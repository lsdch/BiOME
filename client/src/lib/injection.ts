import { InstanceSettings } from '@/api'
import { inject, InjectionKey } from 'vue'

export const settingsKey = Symbol('settings') as InjectionKey<InstanceSettings>

export function injectSettings(): InstanceSettings {
  const settings = inject(settingsKey)
  if (!settings) {
    throw new Error('InstanceSettings not provided')
  }
  return settings
}
