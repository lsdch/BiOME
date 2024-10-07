<template>
  <v-avatar color="white" border="md">
    <v-img :aspect-ratio="1" rounded="circle" v-bind="iconImgProps" :max-width></v-img>
  </v-avatar>
</template>

<script setup lang="ts">
import { useInstanceSettings } from '../settings'
const { iconImgProps } = useInstanceSettings()

defineProps<{ maxWidth?: number }>()
</script>

<style scoped></style>
