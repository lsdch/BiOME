<template>
  <v-menu
    location="bottom end"
    :width="300"
    :min-width="100"
    :close-on-content-click="false"
    :offset="9"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        :color="active ? 'primary' : 'secondary'"
        dark
        v-bind="props"
        icon="mdi-filter-menu"
        variant="plain"
      />
    </template>

    <v-list>
      <v-list-subheader class="text-overline"> Item Filters </v-list-subheader>
      <v-divider />
      <v-list-item>
        <v-switch
          class="px-3"
          label="Owned items"
          v-model="model.ownedItems"
          hide-details
          density="compact"
          color="primary"
          :disabled="user === undefined"
        />
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { User } from '@/api'
import { computed } from 'vue'

const model = defineModel<{
  ownedItems: boolean
}>({ required: true })

defineProps<{ user?: User }>()

const active = computed(() => {
  return Object.values(model.value).find((v) => v) !== undefined
})
</script>

<style scoped></style>
