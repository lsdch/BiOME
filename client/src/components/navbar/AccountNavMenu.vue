<template>
  <v-btn
    v-if="user === undefined"
    :variant="smAndDown ? 'flat' : 'outlined'"
    :prepend-icon="smAndDown ? '' : 'mdi-account-circle'"
    :icon="smAndDown ? 'mdi-account' : undefined"
    :color="smAndDown ? 'primary' : undefined"
    :text="smAndDown ? '' : 'Sign in'"
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
      <v-btn icon="mdi-account-circle" v-bind="props" />
    </template>
    <v-list>
      <v-list-subheader class="mb-3">
        <div class="d-flex align-center">
          <v-icon class="mr-5" v-bind="roleIcon(user.role)" />
          <div>
            <span class="text-overline">
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
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'
import { roleIcon } from '../people/userRole'

const { smAndDown } = useDisplay()

const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const router = useRouter()

const { feedback } = useFeedback()

async function logout() {
  await userStore.logout()
  router.push({ name: 'home' })
  feedback({ type: 'info', message: 'You have been logged out' })
}
</script>

<style scoped></style>
