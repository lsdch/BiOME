<template>
  <v-btn
    v-if="user === undefined"
    :variant="smAndDown ? 'flat' : 'tonal'"
    :prepend-icon="smAndDown ? '' : 'mdi-account'"
    :icon="smAndDown ? 'mdi-account' : undefined"
    :text="smAndDown ? '' : 'Sign in'"
    :to="{
      name: 'login',
      query:
        $router.currentRoute.value.name === 'login'
          ? undefined
          : { redirect: $router.currentRoute.value.path }
    }"
  />
  <v-avatar v-else class="mx-2" variant="outlined" color="blue">
    <v-menu
      location="bottom right"
      :close-on-content-click="false"
      target="#app-bar"
      :width="300"
      :min-width="100"
    >
      <template v-slot:activator="{ props }">
        <v-btn icon="mdi-account" variant="flat" color="light-blue-darken-2" v-bind="props" />
      </template>
      <v-list>
        <v-list-subheader class="mb-3">
          <div class="d-flex align-center">
            <UserRole.Icon class="mr-5" :role="user.role" />
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
  </v-avatar>
</template>

<script setup lang="ts">
import { UserRole } from '@/api'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()

const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const router = useRouter()

const { feedback } = useFeedback()

async function logout() {
  userStore.logout()
  router.push({ name: 'home' })
  feedback({ type: 'info', message: 'You have been logged out' })
}
</script>

<style scoped></style>
