<template>
  <v-container>
    <v-row>
      <v-col lg="6" offset-lg="3">
        <v-card>
          <v-card-text v-if="mode === Mode.Login">
            <LoginForm />
            <div class="d-flex justify-space-between align-center">
              <v-btn
                variant="plain"
                text="I forgot my password"
                :ripple="false"
                class="text-none text-center"
                @click="mode = Mode.PasswordReset"
              />
              <v-btn
                v-if="instance?.allow_contributor_signup"
                size="large"
                color="primary"
                text="Register"
                variant="plain"
                :to="{ name: 'signup' }"
              />
            </div>
          </v-card-text>
          <v-card-text v-if="mode === Mode.PasswordReset">
            <PasswordResetForm />
            <div class="d-flex justify-center">
              <v-btn
                variant="plain"
                text="Back to login page"
                class="text-none"
                :ripple="false"
                @click="mode = Mode.Login"
              ></v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { Ref, ref } from 'vue'
import LoginForm from '@/components/auth/LoginForm.vue'
import PasswordResetForm from '@/components/auth/PasswordResetRequestForm.vue'
import { useInstanceSettings } from '@/components/settings'

enum Mode {
  Login,
  PasswordReset
}

const mode: Ref<Mode> = ref(Mode.Login)

const { instance } = useInstanceSettings()
</script>

<style lang="less"></style>
