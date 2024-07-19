<template>
  <span class="text-overline">{{ site.code }}</span>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'

defineProps<{ site: SiteItem }>()
</script>

<style scoped></style>
