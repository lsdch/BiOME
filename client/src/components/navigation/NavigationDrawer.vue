<template>
  <v-navigation-drawer v-model="drawer" v-bind="$attrs">
    <v-list density="compact" nav open-strategy="single">
      <template v-for="group in navRoutes" :key="group.label">
        <!-- Route -->
        <v-list-item
          v-if="!group.routes"
          v-show="!group.granted || isGranted(group.granted)"
          :prepend-icon="group.icon"
          :title="group.label"
          color="primary"
          :to="$router.resolve(group)"
        />
        <!-- Route group -->
        <v-list-group v-else v-show="!group.granted || isGranted(group.granted)">
          <template v-slot:activator="{ props }">
            <v-list-item
              v-bind="props"
              :prepend-icon="group.icon"
              :title="group.label"
              color="primary"
              :active="group.routes?.find(isRouteActive) !== undefined"
            />
          </template>
          <v-list-item
            v-for="route in group.routes"
            v-show="!route.granted || isGranted(route.granted)"
            :key="route.label"
            :title="route.label"
            link
            slim
            :to="$router.resolve(route)"
            color="primary"
            :active="isRouteActive(route)"
            :prepend-icon="route.icon"
          />
        </v-list-group>
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { RouteDefinition, navRoutes } from '@/router'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const router = useRouter()
const drawer = defineModel<boolean>({ default: true })
function isRouteActive(route: RouteDefinition) {
  return route.name === router.currentRoute.value.name
}

const { isGranted } = useUserStore()
</script>

<style scoped></style>
