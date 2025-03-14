<template>
  <v-navigation-drawer
    id="main-nav-drawer"
    v-model="drawer"
    :location="mobile ? 'top' : 'start'"
    v-bind="$attrs"
  >
    <v-list density="compact" nav open-strategy="single">
      <template v-for="(group, key) in navRoutes" :key>
        <v-divider v-if="isDivider(group)"></v-divider>
        <!-- Route -->
        <v-list-item
          v-else-if="!group.routes"
          v-show="!group.granted || isGranted(group.granted)"
          :prepend-icon="group.icon"
          :title="group.label"
          color="primary"
          :to="$router.resolve(group)"
        />
        <!-- Route group -->
        <v-list-group
          v-else
          v-show="!group.granted || isGranted(group.granted)"
          v-bind="group.groupProps"
        >
          <!-- Route group activator -->
          <template v-slot:activator="{ props }">
            <v-list-item
              v-bind="props"
              :prepend-icon="group.icon"
              :title="group.label"
              color="primary"
              :active="group.routes?.find(isRouteActive) !== undefined"
            />
          </template>
          <!-- Inner route definition -->
          <template v-for="route in group.routes">
            <v-divider v-if="'subgroup' in route" />
            <v-list-subheader v-if="'subgroup' in route">
              {{ route.subgroup }}
            </v-list-subheader>
            <v-divider v-if="'subgroup' in route" />

            <v-list-item
              v-else
              v-show="!route.granted || isGranted(route.granted)"
              :key="route.label"
              :title="route.label"
              link
              slim
              :to="$router.resolve(route)"
              color="primary"
              :active="isRouteActive(route)"
              :prepend-icon="route.icon"
              v-bind="route.itemProps"
            />
          </template>
        </v-list-group>
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { RouteDefinition, navRoutes, isDivider, RouteSubgroup } from '@/router'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'

const { mobile } = useDisplay()

const router = useRouter()
const drawer = defineModel<boolean>({ default: true })
function isRouteActive(route: RouteDefinition | RouteSubgroup) {
  return 'name' in route && route.name === router.currentRoute.value.name
}

const { isGranted } = useUserStore()
</script>

<style lang="scss">
@use 'vuetify';
#main-nav-drawer {
  .v-list-item__spacer {
    width: 18px;
  }

  .v-list-group__items .v-list-item,
  .v-list-subheader {
    padding-inline-start: var(--indent-padding) !important;
  }
}
</style>
