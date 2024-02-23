<template>
  <v-btn
    v-if="user === undefined"
    :variant="smAndDown ? 'flat' : 'outlined'"
    :prepend-icon="smAndDown ? '' : 'mdi-account-circle'"
    :icon="smAndDown ? 'mdi-account' : undefined"
    :color="smAndDown ? 'primary' : undefined"
    :text="smAndDown ? '' : 'Sign in'"
    :to="{ name: 'login' }"
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
      <v-btn icon="mdi-account-circle" v-bind="props" />
    </template>
    <v-list>
      <v-list-subheader class="mb-3">
        <span class="text-overline">
          {{ user.identity.first_name }} {{ user.identity.last_name }}
        </span>
        <br />
        <span class="text-caption">{{ user.role }}</span>
      </v-list-subheader>
      <v-divider />
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
  <v-snackbar v-model="snackbar" multi-line :timeout="3000">
    {{ snackbarText }}

    <template v-slot:actions>
      <v-btn color="primary" variant="text" @click="snackbar = false"> Close </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useDisplay } from 'vuetify'
import { useUserStore } from '@/stores/user'

const { smAndDown } = useDisplay()

const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const snackbar = ref(false)
const snackbarText = 'You are now logged out'
const router = useRouter()

async function logout() {
  await userStore.logout()
  router.push({ name: 'home' })
  snackbar.value = true
}
</script>

<style scoped></style>
