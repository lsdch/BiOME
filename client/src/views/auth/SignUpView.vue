<template>
  <v-card
    variant="flat"
    :title="registrationDone ? 'New account registered' : 'Apply for a contributor account'"
  >
    <template #prepend>
      <v-icon
        :icon="registrationDone ? 'mdi-check-bold' : 'mdi-account-plus'"
        :color="registrationDone ? 'green' : 'primary'"
      />
    </template>
    <v-card-text v-if="registrationDone">
      An email was sent to your address with a link to activate your account. Please check your
      inbox.
    </v-card-text>
    <v-card-text v-else>
      <v-container>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field
              label="First name"
              v-model="model.first_name"
              v-bind="field('first_name')"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field label="Last name" v-model="model.last_name" v-bind="field('last_name')" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field
              label="E-mail address"
              v-model="model.email"
              prepend-inner-icon="mdi-at"
              v-bind="field('email')"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field
              label="Institution"
              v-model="model.institution"
              placeholder="(Optional) The institution your attached to, e.g. your lab or association"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-textarea
              label="Motivations"
              v-model="model.motive"
              v-bind="field('motive')"
              placeholder="(Optional) Short explanation on how you want to contribute"
              persistent-placeholder
            />
          </v-col>
        </v-row>
        <v-btn color="primary" block text="Submit" @click="submit" />
      </v-container>
      <div class="d-flex justify-center mt-3">
        <v-btn
          class="text-none"
          :to="{ name: 'login' }"
          variant="plain"
          :ripple="false"
          text="Back to login page"
        />
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { $PendingUserRequestInput, AccountService, PendingUserRequestInput } from '@/api'
import { useSchema } from '@/components/toolkit/forms/schema'
import { accountRoutes } from '@/router/routes'
import { ref } from 'vue'

const registrationDone = ref(false)
const model = ref<PendingUserRequestInput>({
  email: '',
  motive: undefined,
  first_name: '',
  last_name: '',
  institution: undefined
})

const { field, errorHandler } = useSchema($PendingUserRequestInput)

function submit() {
  AccountService.register({
    body: {
      data: model.value,
      verification_path: accountRoutes.verifyEmail.path
    }
  })
    .then(errorHandler)
    .then(() => (registrationDone.value = true))
}
</script>

<style scoped></style>
