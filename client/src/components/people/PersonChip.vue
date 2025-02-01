<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip
        :text="short ? `${person.first_name[0]}. ${person.last_name}` : person.full_name"
        v-bind="{ ...props, ...$attrs }"
      />
    </template>
    <v-card
      :title="person.full_name"
      :subtitle="person.alias"
      class="bg-surface-light small-card-title"
      density="compact"
    >
      <template #prepend>
        <div class="d-flex flex-column align-center mr-2">
          <UserRole.Icon :role="person.role" size="small" />
        </div>
      </template>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { PersonInner, UserRole } from '@/api'

defineProps<{ person: PersonInner; short?: boolean }>()
</script>

<style scoped lang="scss"></style>
