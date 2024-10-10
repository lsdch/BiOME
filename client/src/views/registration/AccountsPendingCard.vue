<template>
  <v-card hover @click="toggleActive()" :ripple="false">
    <template #append>
      <v-btn color="" variant="plain" :icon="active ? 'mdi-chevron-up' : 'mdi-chevron-down'" />
    </template>
    <template #title> {{ item.identity.first_name }} {{ item.identity.last_name }} </template>
    <template #subtitle>
      {{ item.email }}
      <v-icon
        size="small"
        v-bind="
          item.email_verified
            ? {
                icon: 'mdi-check-circle',
                color: 'success'
              }
            : {
                icon: 'mdi-clock',
                color: 'warning'
              }
        "
      />
    </template>
    <div v-show="active">
      <v-btn
        class="ml-3"
        color="primary"
        text="Invite"
        prepend-icon="mdi-account-plus"
        variant="plain"
        @click.stop="emit('invite', item)"
      />
      <v-list>
        <v-list-item title="Institution" :subtitle="item.institution"></v-list-item>
        <v-list-item title="Motivations" :subtitle="item.motive"></v-list-item>
      </v-list>
    </div>
    <template #actions>
      <ItemDateChip :date="item.created_on" icon="mdi-calendar-clock" />
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { PendingUserRequest } from '@/api'
import ItemDateChip from '@/components/toolkit/ItemDateChip.vue'
import { useToggle } from '@vueuse/core'

defineProps<{ item: PendingUserRequest }>()

const [active, toggleActive] = useToggle(false)

const emit = defineEmits<{
  invite: [item: PendingUserRequest]
}>()
</script>

<style scoped></style>
