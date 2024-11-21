<template>
  <v-bottom-navigation v-if="mobile" mandatory>
    <v-btn
      v-for="{ title, icon, category } in subroutes"
      :value="category"
      :key="category"
      :text="title"
      :prepend-icon="icon"
      :to="resolveSubroute(category)"
      color="primary"
    />
  </v-bottom-navigation>
  <v-navigation-drawer v-else permanent :width="200">
    <v-list nav>
      <v-list-subheader>Settings</v-list-subheader>
      <v-list-item
        v-for="{ title, icon, category } in subroutes"
        :key="category"
        :title
        color="primary"
        :prepend-icon="icon"
        :to="resolveSubroute(category)"
      />
    </v-list>
  </v-navigation-drawer>
  <div class="bg-surface d-flex fill-height flex-column">
    <v-container :max-width="1200">
      <Suspense>
        <template #default>
          <component :is="component" />
        </template>
        <template #fallback>
          <v-skeleton-loader type="article, article, article" />
        </template>
      </Suspense>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import EmailSettings from '@/components/settings/EmailSettings.vue'
import InstanceSettings from '@/components/settings/InstanceSettings.vue'
import SecuritySettings from '@/components/settings/SecuritySettings.vue'
import { computed } from 'vue'
import { useDisplay } from 'vuetify'

import NotFound from '@/components/navigation/NotFound.vue'
import routes from '@/router/routes'
import { useRouter } from 'vue-router'

const { mobile } = useDisplay()

type SettingsTab = 'instance' | 'security' | 'email'

const { resolve } = useRouter()

const props = defineProps<{ category: SettingsTab }>()
function resolveSubroute(category: string) {
  return resolve({ name: routes.settings.name, params: { category } })
}
const subroutes = [
  { title: 'Instance', category: 'instance', icon: 'mdi-application-settings-outline' },
  { title: 'Security', category: 'security', icon: 'mdi-security' },
  { title: 'E-mailing', category: 'email', icon: 'mdi-email' }
]

const component = computed(() => {
  switch (props.category) {
    case 'instance':
      return InstanceSettings
    case 'security':
      return SecuritySettings
    case 'email':
      return EmailSettings
    default:
      return NotFound
  }
})
</script>

<style scoped></style>
