<template>
  <v-container>
    <HomeLinkTitle />

    <v-row>
      <v-col cols="12" lg="6" offset-lg="3" class="text-center">
        <div v-if="result === undefined">
          Confirming your account, please wait...
          <v-progress-circular indeterminate color="primary" />
        </div>
        <v-alert
          v-else-if="result === EmailConfirmationError.AlreadyVerified"
          type="warning"
          variant="outlined"
          text="This account was already activated."
        />
        <v-alert
          v-else-if="result === EmailConfirmationError.InvalidToken"
          type="warning"
          variant="outlined"
          text="Invalid confirmation token."
        />
        <v-alert v-else-if="result === 'ServerError'" text="An error occurred on the server." />
        <v-alert
          v-else-if="result === 'ConfirmationSuccess'"
          variant="outlined"
          type="success"
          text="Your account has been activated! You have been logged in. Redirecting to home page..."
        ></v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ApiError, AccountService, EmailConfirmationError } from '@/api'
import HomeLinkTitle from '@/components/navigation/HomeLinkTitle.vue'
import { Ref } from 'vue'
import { ref } from 'vue'
import { onMounted } from 'vue'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const token = computed(() => route.query.token?.toString() ?? '')

type Response = 'ConfirmationSuccess' | EmailConfirmationError | 'ServerError'
const result: Ref<undefined | Response> = ref(undefined)

onMounted(() => {
  AccountService.confirmEmail({ token: token.value })
    .then(() => {
      result.value = 'ConfirmationSuccess'
      setTimeout(() => {
        router.push({ name: 'home' })
      }, 3000)
    })
    .catch((err: ApiError) => {
      switch (err.status) {
        case 400:
          result.value = err.body
          break

        default:
          result.value = 'ServerError'
          break
      }
    })
})
</script>

<style scoped></style>
