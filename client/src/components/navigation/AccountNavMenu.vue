<template>
  <v-btn
    v-if="user === undefined"
    variant="outlined"
    prepend-icon="mdi-account-circle"
    text="Sign in"
    :to="{
      name: 'login',
      query:
        $router.currentRoute.value.name === 'login'
          ? undefined
          : { redirect: $router.currentRoute.value.path }
    }"
  />
  <v-menu
    v-else
    location="bottom right"
    :close-on-content-click="false"
    target="#app-bar"
    :width="300"
    :min-width="100"
  >
    <template v-slot:activator="{ props }">
      <v-btn variant="outlined" v-bind="props" class="text-lg-body-1">
        <template #prepend>
          <UserRole.Icon :role="user.role" />
        </template>
        {{
          user.identity.first_name
            .split(/[ -]/)
            .map((w) => w[0])
            .join('')
        }}
        {{ user.identity.last_name }}
      </v-btn>
    </template>
    <v-list>
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
      <v-divider />
      <v-list-item>
        <v-btn
          prepend-icon="mdi-account"
          variant="plain"
          density="compact"
          text="Account"
          :to="{ name: 'account' }"
        />
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
</template>

<script setup lang="ts">
import { UserRole } from '@/api'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { watch } from 'vue'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'

const userStore = useUserStore()
const { user, usePrivilege } = storeToRefs(userStore)

import { useTheme } from 'vuetify'
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
