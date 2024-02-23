<template>
  <v-navigation-drawer v-model="drawer">
    <v-list density="compact" nav>
      <template v-for="group in routeGroups" :key="group.label">
        <v-list-item
          v-if="!group.routes"
          :prepend-icon="group.icon"
          :title="group.label"
          color="primary"
          :to="$router.resolve(group)"
        />
        <v-list-group v-else>
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
import { useRouter } from 'vue-router'
import { RouteDefinition, routeGroups } from './router'
const drawer = defineModel<boolean>({ default: true })
const router = useRouter()
function isRouteActive(route: RouteDefinition) {
  return route.name === router.currentRoute.value.name
}
</script>

<style scoped></style>
