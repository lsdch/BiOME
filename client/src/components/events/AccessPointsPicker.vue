<template>
  <v-combobox :items :loading chips closable-chips multiple>
    <template #chip="{ item, props }">
      <v-chip
        v-bind="props"
        :text="item.value"
        :color="items.includes(item.value) ? 'primary' : 'success'"
      />
    </template>
  </v-combobox>
</template>

<script setup lang="ts">
import { SamplingService } from '@/api'
import { useToggle } from '@vueuse/core'
import { onMounted, ref } from 'vue'

const items = ref<string[]>([])

const [loading, toggleLoading] = useToggle(true)

async function fetch() {
  toggleLoading(true)
  return SamplingService.getAccessPoints({ throwOnError: true })
    .then(({ data }) => data)
    .catch((err) => {
      console.error('Failed to retrieve list of access points', err)
      return []
    })
    .finally(() => toggleLoading(false))
}

onMounted(async () => (items.value = await fetch()))
</script>

<style scoped lang="scss"></style>
