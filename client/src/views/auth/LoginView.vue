<template>
  <v-container class="fill-height">
    <v-row :align="smAndDown ? 'baseline' : 'center'">
      <v-col lg="6" offset-lg="3">
        <div class="w-100 d-flex align-center flex-column mb-5">
          <v-btn
            class="plain-link text-h3 font-weight-bold"
            :to="{ name: 'home' }"
            :text="app_name"
            variant="text"
            size="x-large"
            :ripple="false"
            plain
          />
        </div>
        <v-card :variant="smAndDown ? 'flat' : 'elevated'">
          <v-card-text v-if="mode === 'Login'">
            <LoginForm></LoginForm>
            <div class="d-flex justify-space-between align-center">
              <v-btn
                variant="plain"
                text="I forgot my password"
                class="text-none text-center"
                @click="mode = 'PasswordReset'"
              ></v-btn>
              <v-btn
                size="large"
                color="primary"
                text="Register"
                variant="plain"
                :to="{ name: 'signup' }"
              />
            </div>
          </v-card-text>
          <v-card-text v-if="mode === 'PasswordReset'">
            <PasswordResetForm />
            <div class="d-flex justify-center">
              <v-btn
                variant="plain"
                text="Back to login page"
                class="text-none"
                @click="mode = 'Login'"
              ></v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { Ref } from 'vue'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'
import LoginForm from '@/components/auth/LoginForm.vue'
import PasswordResetForm from '@/components/auth/PasswordResetForm.vue'

type Mode = 'Login' | 'PasswordReset'

const { smAndDown } = useDisplay()

const app_name = import.meta.env.VITE_APP_NAME

const mode: Ref<Mode> = ref('Login')
</script>

<style lang="less">
.plain-link .v-btn__overlay {
  display: none;
}
</style>
