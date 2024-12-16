<template>
  <v-app>
    <v-app-bar id="app-bar" color="primary" v-if="!$route.meta.hideNavbar" density="compact">
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-app-bar-title :to="$router.resolve({ name: 'home' })">
        <v-btn
          class="app-title opacity-100"
          variant="plain"
          :ripple="false"
          :to="{ name: 'home' }"
          :text="xs ? undefined : settings.name"
        >
          <template #prepend>
            <AppIcon :size="30" />
          </template>
        </v-btn>
      </v-app-bar-title>
      <v-spacer />
      <SettingsMenu />
      <AccountNavMenu />
    </v-app-bar>
    <NavigationDrawer v-model="drawer" :temporary="drawerTemporary || smAndDown" />
    <v-main id="main" class="bg-main">
      <v-progress-linear v-show="loading" :color="colors.orange.base" indeterminate />
      <RouterView :key="$route.fullPath" v-slot="{ Component }">
        <Suspense>
          <component :is="Component" />
          <template #fallback>
            <v-container class="fill-height d-flex align-center justify-center">
              <v-card class="d-flex align-center justify-center" min-height="50%" :min-width="500">
                <v-progress-circular indeterminate size="large"></v-progress-circular>
              </v-card>
            </v-container>
          </template>
        </Suspense>
      </RouterView>
    </v-main>
    <ErrorSnackbar v-model="snackbar.open" :title="snackbar.title" :errors="snackbar.errors" />
    <FeedbackSnackbar />
    <ConfirmDialog
      :model-value="isRevealed"
      v-bind="content"
      @confirm="confirm()"
      @cancel="cancel"
    />
  </v-app>
</template>

<script setup lang="ts">
import NavigationDrawer from '@/components/navigation/NavigationDrawer.vue'
import { nextTick, ref } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import colors from 'vuetify/util/colors'

import AccountNavMenu from '@/components/navbar/AccountNavMenu.vue'
import SettingsMenu from '@/components/navbar/SettingsMenu.vue'

import { client, ErrorDetail, InstanceSettings } from '@/api'
import { useDisplay } from 'vuetify'
import AppIcon from './components/icons/AppIcon.vue'
import ConfirmDialog from './components/toolkit/ui/ConfirmDialog.vue'
import ErrorSnackbar from './components/toolkit/ui/ErrorSnackbar.vue'
import FeedbackSnackbar from './components/toolkit/ui/FeedbackSnackbar.vue'
import { useAppConfirmDialog } from './composables/confirm_dialog'

const loading = ref(false)

const { smAndDown, xs } = useDisplay()

const drawer = ref(!smAndDown.value)
const drawerTemporary = ref<boolean>()

// Navigation
const router = useRouter()
router.beforeEach(async (to) => {
  drawerTemporary.value = to.meta.drawer?.temporary
  await nextTick()
  loading.value = true
})
router.afterEach((to) => {
  if (to.name === 'api-docs') {
    drawer.value = false
  }
  loading.value = false
})

defineProps<{ settings: InstanceSettings }>()

// Confirm dialog
const { isRevealed, confirm, cancel, content } = useAppConfirmDialog()

// Access control feedback
const snackbar = ref<{ open: boolean; title: string; errors: ErrorDetail[] }>({
  open: false,
  title: '',
  errors: []
})

client.interceptors.response.use(async (response) => {
  if (response.status === 403) {
    snackbar.value.title = 'Access denied'
    snackbar.value.open = true
  }
  return response
})
</script>

<style>
#main {
  background: initial;
  background-image: radial-gradient(#64646433 1px, transparent 0px);
  background-size: 25px 25px;
  background-position: -10px -10px;
}

.app-title {
  color: white;
  font-weight: bold;
  font-size: larger;
  text-transform: none;
  font-family: Verdana, Geneva, Tahoma, sans-serif;
}

.noTransition {
  transition-duration: 0s !important;
}
</style>
