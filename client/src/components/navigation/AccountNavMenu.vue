<template>
  <v-menu
    location="bottom right"
    :close-on-content-click="false"
    target="#app-bar"
    :width="300"
    :min-width="100"
  >
    <template v-slot:activator="{ props }">
      <v-btn icon v-bind="props" :rounded="100">
        <v-avatar
          color="surface-variant"
          :badge="
            user
              ? {
                  color: RoleIcon[usePrivilege ?? user.role].color,
                  location: 'bottom end'
                }
              : undefined
          "
        >
          <span v-if="!!user" class="on-surface font-weight-bold">
            {{ userInitials }}
          </span>
          <v-icon v-else class="on-surface" icon="mdi-account" />
        </v-avatar>
      </v-btn>
    </template>
    <v-list>
      <template v-if="!!user">
        <v-list-subheader class="mb-3">
          <div class="d-flex align-center">
            <UserRole.Icon class="mr-5" :role="user.role" />
            <div>
              <span class="font-weight-bold">
                {{ user.identity.full_name }}
              </span>
              <br />
              <span class="text-caption">
                {{ user.role }}
              </span>
            </div>
          </div>
        </v-list-subheader>
        <v-divider />
        <v-list-item>
          <v-select
            v-model="usePrivilege"
            label="Use role privileges"
            :items="UserRole.upTo(user.role, false)"
            hide-details
            variant="solo-filled"
          />
        </v-list-item>
        <v-divider />
        <v-list-item prepend-icon="mdi-account" title="Account" :to="{ name: 'account' }">
        </v-list-item>
        <v-list-item prepend-icon="mdi-power" title="Logout" @click="logout()"> </v-list-item>
      </template>
      <template v-else>
        <v-list-item
          title="Sign in"
          prepend-icon="mdi-login"
          :to="{
            name: 'login',
            query:
              $router.currentRoute.value.name === 'login'
                ? undefined
                : { redirect: $router.currentRoute.value.path }
          }"
        ></v-list-item>
      </template>
      <v-divider></v-divider>
      <v-list-item>
        <v-switch
          class="px-3"
          label="Dark theme"
          :model-value="theme.name.value"
          @update:model-value="(v) => theme.change(v!)"
          false-value="light"
          true-value="dark"
          color="purple"
          :indeterminate="false"
          hide-details
        />
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { UserRole } from '@/api'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { computed, watch } from 'vue'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const { user, usePrivilege } = storeToRefs(userStore)

const userInitials = computed(() => {
  if (!user.value) return ''
  return user.value.identity.first_name[0] + user.value.identity.last_name[0]
})

import { useTheme } from 'vuetify'
import { RoleIcon } from '../icons/UserRoleIcon'
const theme = useTheme()

watch(theme.name, () => {
  localStorage.setItem('app-theme', theme.name.value)
})

const { feedback } = useFeedback()
const router = useRouter()
async function logout() {
  userStore.logout()
  router.push({ name: 'home' })
  feedback({ type: 'info', message: 'You have been logged out' })
}
</script>

<style scoped></style>
