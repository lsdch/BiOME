<template>
  <v-container>
    <HomeLinkTitle />
    <v-row>
      <v-col cols="12" lg="8" offset-lg="2">
        <v-card :variant="smAndDown ? 'flat' : 'elevated'">
          <v-card-title primary-title> Sign up for an account </v-card-title>
          <v-card-text>
            <SignUpForm />
            <div class="d-flex justify-center mt-3">
              <v-btn class="text-none" :to="{ name: 'login' }" variant="plain" :ripple="false">
                Back to login page
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import SignUpForm from '@/components/auth/SignUpForm.vue'
import HomeLinkTitle from '@/components/navigation/HomeLinkTitle.vue'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()
</script>

<style scoped></style>
