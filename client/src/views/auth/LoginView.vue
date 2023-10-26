<template>
  <v-container class="fill-height">
    <v-row :align="smAndDown ? 'baseline' : 'center'">
      <v-col lg="6" offset-lg="3">
        <HomeLinkTitle />
        <v-card :variant="smAndDown ? 'flat' : 'elevated'">
          <v-card-text v-if="mode === 'Login'">
            <LoginForm></LoginForm>
            <div class="d-flex justify-space-between align-center">
              <v-btn
                variant="plain"
                text="I forgot my password"
                :ripple="false"
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
                :ripple="false"
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
import { Ref, ref } from 'vue'
import { useDisplay } from 'vuetify'
import LoginForm from '@/components/auth/LoginForm.vue'
import PasswordResetForm from '@/components/auth/PasswordResetForm.vue'
import HomeLinkTitle from '@/components/navigation/HomeLinkTitle.vue'

type Mode = 'Login' | 'PasswordReset'

const { smAndDown } = useDisplay()

const mode: Ref<Mode> = ref('Login')
</script>

<style lang="less"></style>
