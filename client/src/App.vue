<template>
  <v-app>
    <v-app-bar id="app-bar" color="primary" v-if="!$route.meta.hideNavbar">
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-app-bar-title :to="$router.resolve({ name: 'home' })">
        <v-btn
          class="app-title"
          variant="plain"
          :ripple="false"
          :to="{ name: 'home' }"
          :text="APP_TITLE"
        />
      </v-app-bar-title>
      <v-spacer />
      <SettingsMenu />
      <AccountNavMenu />
    </v-app-bar>
    <NavigationDrawer v-model="drawer" />
    <v-main>
      <RouterView />
    </v-main>
    <ErrorSnackbar v-model="snackbar.open" :title="snackbar.title" :errors="snackbar.errors" />
  </v-app>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import NavigationDrawer from '@/components/navigation/NavigationDrawer.vue'

import AccountNavMenu from '@/components/navbar/AccountNavMenu.vue'
import SettingsMenu from '@/components/navbar/SettingsMenu.vue'

import { ErrorDetail, OpenAPI } from './api'
import ErrorSnackbar from './components/toolkit/ui/ErrorSnackbar.vue'

const snackbar = ref<{ open: boolean; title: string; errors: ErrorDetail[] }>({
  open: false,
  title: '',
  errors: []
})

OpenAPI.interceptors.response.use(async (response) => {
  if (response.status === 401) {
    const body = await response.json()
    snackbar.value.title = 'Access denied'
    snackbar.value.errors = body.errors
    snackbar.value.open = true
  }
  return response
})

const APP_TITLE = import.meta.env.VITE_APP_NAME
const drawer = ref(true)

const router = useRouter()
router.afterEach((to) => {
  if (to.name === 'api-docs') {
    drawer.value = false
  }
})
</script>

<style>
.app-title {
  color: white;
  font-weight: bold;
  font-size: larger;
  text-transform: none;
  font-family: Verdana, Geneva, Tahoma, sans-serif;
}
</style>
