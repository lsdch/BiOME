<template>
  <v-text-field
    v-model="model"
    v-bind="$attrs"
    :type="show ? 'text' : 'password'"
    :append-inner-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
    @click:append-inner="show = !show"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'

const model = defineModel<string>()
const show = ref(false)
</script>

<style scoped></style>
