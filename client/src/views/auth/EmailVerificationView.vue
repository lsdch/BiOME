<template>
  <v-container>
    <v-row>
      <v-col cols="12" lg="6" offset-lg="3" class="text-center">
        <v-card v-if="result === undefined">
          <template #text>
            <p>Attempting e-mail verification...</p>
            <v-progress-circular class="my-5" indeterminate color="primary" />
          </template>
        </v-card>
        <v-alert v-else-if="result === true" type="success" variant="tonal">
          <p class="text-h6">Your e-mail address has been verified!</p>
          <p class="my-5">
            Once an administrator has reviewed your request,<br />
            you will receive another e-mail with a link to register an account.
          </p>
          <v-btn color="primary" variant="plain" :to="{ name: 'home' }">Back to home page</v-btn>
        </v-alert>
        <v-alert v-else variant="tonal" type="error">
          <p>E-mail verification failed.</p>
          {{ result.detail }}
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { AccountService, ErrorModel } from '@/api'
import { useRouteQuery } from '@vueuse/router'
import { computed, onMounted, Ref, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

const token = useRouteQuery<string>('token', undefined)

const result = ref<true | ErrorModel>()

onMounted(async () => {
  const { data, error } = await AccountService.confirmEmail({ query: { token: token.value } })
  if (error) {
    result.value = error
  } else {
    result.value = true
  }
})
</script>

<style scoped></style>
