<template>
  <v-card
    title="Event"
    :subtitle="siteCode ? 'From site' : 'Waiting for site definition'"
    class="fill-height small-card-title"
    prepend-icon="mdi-calendar"
  >
    <template #append v-if="siteCode">
      <v-btn text="New event" />
    </template>
    <v-card-text>
      <SiteEventPicker v-if="siteCode" v-model="event" :site-code clearable />
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { Event } from '@/api'
import { watchEffect } from 'vue'
import SiteEventPicker from './SiteEventPicker.vue'

const event = defineModel<Event>()
const props = defineProps<{ siteCode?: string }>()
watchEffect(() => {
  if (!props.siteCode) event.value = undefined
})
</script>

<style scoped lang="scss"></style>
