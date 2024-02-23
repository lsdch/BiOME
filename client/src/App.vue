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
  </v-app>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import NavigationDrawer from './NavigationDrawer.vue'

import AccountNavMenu from '@/components/navbar/AccountNavMenu.vue'
import SettingsMenu from '@/components/navbar/SettingsMenu.vue'

const APP_TITLE = import.meta.env.VITE_APP_NAME
const drawer = ref(true)
</script>

<style lang="less">
.app-title {
  color: white;
  font-weight: bold;
  font-size: larger;
  text-transform: none;
  font-family: Verdana, Geneva, Tahoma, sans-serif;
}
</style>
