<template>
  <v-avatar color="white" border="md" :size>
    <v-img :aspect-ratio="1" rounded="circle" v-bind="iconImgProps" :size></v-img>
  </v-avatar>
</template>

<script setup lang="ts">
import { useInstanceSettings } from '../settings'
const { iconImgProps } = useInstanceSettings()

defineProps<{ size?: number }>()
</script>

<style scoped></style>
