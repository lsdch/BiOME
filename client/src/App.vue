<template>
  <v-app>
    <v-app-bar id="app-bar" color="primary" v-if="!$route.meta.hideNavbar" density="compact">
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-app-bar-title :to="{ name: 'home' }">
        <v-btn
          class="app-title opacity-100"
          variant="plain"
          :ripple="false"
          :to="{ name: 'home' }"
          :text="xs ? undefined : settings.instance.value?.name"
          :loading="settings.isPending.value"
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
    <NavigationDrawer v-model="drawer" :temporary="lgAndDown && (drawerTemporary || smAndDown)" />
    <v-main id="main" class="bg-main overflow-y-auto" max-height="100vh">
      <v-progress-linear v-show="loading" :color="colors.orange.base" indeterminate />
      <RouterView :key="$route.fullPath" v-slot="{ Component }">
        <Suspense>
          <div id="router-view-suspense-container" class="fill-height">
            <component :is="Component" />
          </div>
          <template #fallback>
            <v-card class="d-flex align-center justify-center fill-height w-100" flat>
              <v-progress-circular indeterminate size="large" color="primary"></v-progress-circular>
            </v-card>
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
    <v-bottom-sheet
      :model-value="true"
      v-if="!cookiesAccepted"
      persistent
      :scrim="false"
      class="rounded-0"
      no-click-animation
    >
      <v-card flat :rounded="0">
        <v-alert icon="mdi-information" rounded="0">
          This website uses cookies to provide a better user experience. We <b>do not track</b> or
          share your activity.
          <template #append>
            <v-btn color="success" text="OK" @click="cookiesAccepted = true" />
          </template>
        </v-alert>
      </v-card>
    </v-bottom-sheet>
  </v-app>
</template>

<script setup lang="ts">
import NavigationDrawer from '@/components/navigation/NavigationDrawer.vue'
import { nextTick, ref } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import colors from 'vuetify/util/colors'

import AccountNavMenu from '@/components/navbar/AccountNavMenu.vue'
import SettingsMenu from '@/components/navbar/SettingsMenu.vue'

import { client } from '@/api/gen/client.gen'
import { ErrorDetail } from '@/api'
import { useLocalStorage } from '@vueuse/core'
import { useDisplay } from 'vuetify'
import AppIcon from './components/icons/AppIcon.vue'
import { useInstanceSettings } from './components/settings'
import ConfirmDialog from './components/toolkit/ui/ConfirmDialog.vue'
import ErrorSnackbar from './components/toolkit/ui/ErrorSnackbar.vue'
import FeedbackSnackbar from './components/toolkit/ui/FeedbackSnackbar.vue'
import { useAppConfirmDialog } from './composables/confirm_dialog'

const loading = ref(false)

const { lgAndDown, smAndDown, xs } = useDisplay()

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

const settings = useInstanceSettings()

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

const cookiesAccepted = useLocalStorage('cookies-accepted', false)
</script>

<style>
#main {
  /* background: initial;
  background-image: radial-gradient(#64646433 1px, transparent 0px);
  background-size: 25px 25px;
  background-position: -10px -10px; */
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
