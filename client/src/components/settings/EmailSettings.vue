<template>
  <v-form>
    <v-text-field label="SMTP Host" />
    <VNumberInput label="SMTP Port" :min="1" />
    <v-text-field label="User" />
    <v-text-field label="Password" type="password" />
  </v-form>
</template>

<script setup lang="ts"></script>

<style scoped></style>
