<template>
  <v-text-field
    v-model="model"
    v-bind="$attrs"
    :type="show ? 'text' : 'password'"
    :append-inner-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
    @click:append-inner="show = !show"
  >
    <template v-for="(name, index) of slots" #[name]="scope" :key="index">
      <slot :name="name" v-bind="{ scope }" />
    </template>
  </v-text-field>
</template>

<script setup lang="ts">
import { ref, useSlots } from 'vue'

const model = defineModel<string>({ required: true })
const show = ref(false)

const slots = Object.keys(useSlots()) as {}
</script>

<style scoped></style>
