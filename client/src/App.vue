<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { ref } from 'vue'

import { routeGroups } from './router'

const drawer = ref(true)

const router = useRouter()
</script>

<template>
  <v-app>
    <v-app-bar color="primary">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-app-bar-title :to="$router.resolve({ name: 'home' })">
        <v-btn
          class="app-title"
          variant="plain"
          :ripple="false"
          :to="{ name: 'home' }"
          text="DarCo"
        />
      </v-app-bar-title>
    </v-app-bar>

    <!-- <v-navigation-drawer :rail="drawer"> -->
    <v-navigation-drawer v-model="drawer">
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
              <!-- <v-menu v-if="drawer" location="end" activator="parent" offset="10">
                <v-list density="compact" class="elevation-1">
                  <v-list-item
                    v-for="route in group.routes"
                    :key="route.name"
                    :title="route.label"
                    :to="route"
                    color="primary"
                  >
                  </v-list-item>
                </v-list>
              </v-menu> -->
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
