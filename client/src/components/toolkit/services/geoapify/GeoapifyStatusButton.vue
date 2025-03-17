<template>
  <v-menu v-if="status" location="top center">
    <template #activator="{ props }">
      <v-btn v-bind="{ ...props, ...Geoapify.Status.props(status) }" variant="tonal" size="small">
      </v-btn>
    </template>
    <v-card
      class="small-card-title"
      :prepend-icon="status.available ? 'mdi-check' : 'mdi-alert'"
      :title="status?.available ? 'Geoapify API ready' : 'Geoapify API Unavailable'"
      :subtitle
    >
      <v-card-text v-if="!status.has_api_key">
        <v-btn
          v-if="isGranted('Admin')"
          block
          :rounded="0"
          :to="{ name: 'app-settings', params: { category: 'services' } }"
          text="Configure"
        />
        <div v-else class="text-caption">
          Geoapify API was not setup by the administrators.<br />
          Automatic reverse geocoding is disabled.
        </div>
      </v-card-text>
      <v-progress-linear
        v-if="status.has_api_key"
        :max="status.limit"
        :model-value="status.requests"
        color="primary"
      />
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Geoapify } from '@/api'
import { useGeoapify } from '@/stores/geoapify'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'

const { status } = storeToRefs(useGeoapify())
const { isGranted } = useUserStore()

const subtitle = computed(() => {
  if (!status.value) return undefined
  if (!status.value.has_api_key) return 'No API key configured'
  return `${status.value.requests} / ${status.value.limit} daily requests`
})
</script>

<style scoped lang="scss"></style>
