<template>
  <v-container class="fill-height">
    <v-row :align="smAndDown ? 'baseline' : 'center'">
      <v-col lg="6" offset-lg="3">
        <HomeLinkTitle />
        <v-card :variant="smAndDown ? 'flat' : 'elevated'"> </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()
</script>

<style scoped></style>
