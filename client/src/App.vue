<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { ref } from 'vue'

import { routeGroups } from './router'
import { useDisplay, useTheme } from 'vuetify'
import { useUserStore } from './stores/user'
import { storeToRefs } from 'pinia'

const theme = useTheme()

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

const drawer = ref(true)

const router = useRouter()
const { smAndDown } = useDisplay()

const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const APP_TITLE = import.meta.env.VITE_APP_NAME

const snackbar = ref(false)
const snackbarText = 'You are now logged out'

async function logout() {
  await userStore.logout()
  router.push({ name: 'home' })
  snackbar.value = true
}
</script>

<template>
  <v-app>
    <v-app-bar color="primary" v-if="!$route.meta.hideNavbar">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
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
      <v-menu location="bottom" :close-on-content-click="false">
        <template v-slot:activator="{ props }">
          <v-btn icon="mdi-cog" v-bind="props"></v-btn>
        </template>
        <v-list>
          <v-list-item>
            <v-switch
              label="Dark theme"
              v-model="theme.global.current.value.dark"
              @click="toggleTheme"
              color="purple"
            />
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn
        v-if="user === undefined"
        :variant="smAndDown ? 'flat' : 'outlined'"
        :prepend-icon="smAndDown ? '' : 'mdi-account-circle'"
        :icon="smAndDown ? 'mdi-account' : undefined"
        :color="smAndDown ? 'primary' : undefined"
        :text="smAndDown ? '' : 'Sign in'"
        :to="{ name: 'login' }"
      />
      <v-menu v-else location="bottom" :close-on-content-click="false">
        <template v-slot:activator="{ props }">
          <v-btn icon="mdi-account-circle" v-bind="props"></v-btn>
        </template>
        <v-list>
          <v-list-subheader class="mb-3">
            <span class="text-overline">
              {{ user.identity.first_name }} {{ user.identity.last_name }}
            </span>
            <br />
            <span class="text-caption">{{ user.role }}</span>
          </v-list-subheader>
          <v-divider></v-divider>
          <v-list-item>
            <v-btn prepend-icon="mdi-account" variant="plain" density="compact" text="Account" />
          </v-list-item>
          <v-list-item>
            <v-btn
              prepend-icon="mdi-power"
              variant="plain"
              density="compact"
              text="Logout"
              @click="logout()"
            />
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>

    <!-- <v-navigation-drawer :rail="drawer"> -->
    <v-navigation-drawer v-if="!$route.meta.hideNavbar" v-model="drawer">
      <v-list density="compact" nav>
        <div v-for="group in routeGroups" :key="group.name">
          <v-list-item
            v-if="group.routes.length === 1"
            :prepend-icon="group.icon"
            :title="group.routes[0].label"
            :to="group.routes[0]"
            color="primary"
          />
          <v-list-group v-else>
            <template v-slot:activator="{ props }">
              <v-list-item
                v-bind="props"
                :prepend-icon="group.icon"
                :key="group.name"
                :title="group.name"
                color="primary"
                :active="
                  group.routes.find((route) => route.name === router.currentRoute.value.name) !==
                  undefined
                "
              />
            </template>

            <v-list-item
              v-for="route in group.routes"
              :key="route.name"
              :title="route.label"
              :to="route"
              color="primary"
              class="nav-section"
            />
          </v-list-group>
        </div>
        <!-- :to="route" -->
      </v-list>
    </v-navigation-drawer>

    <v-snackbar v-model="snackbar" multi-line :timeout="3000">
      {{ snackbarText }}

      <template v-slot:actions>
        <v-btn color="primary" variant="text" @click="snackbar = false"> Close </v-btn>
      </template>
    </v-snackbar>

    <v-main>
      <RouterView />
    </v-main>
  </v-app>
</template>

<style scoped>
.app-title {
  color: white;
  font-weight: bold;
  font-size: larger;
  text-transform: none;
  font-family: Verdana, Geneva, Tahoma, sans-serif;
}

.v-list-item.nav-section {
  min-height: 30px;
}
</style>
