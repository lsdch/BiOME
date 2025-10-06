<template>
  <v-card hover @click="toggleActive()" :ripple="false">
    <template #prepend>
      <v-checkbox v-model="selected" density="compact" hide-details @click.stop />
    </template>
    <template #append>
      <v-btn color="" variant="plain" :icon="active ? 'mdi-chevron-up' : 'mdi-chevron-down'" />
    </template>
    <template #title> {{ item.full_name }} </template>
    <template #subtitle>
      {{ item.email }}
      <v-tooltip>
        <template #activator="{ props }">
          <v-icon
            size="small"
            v-bind="
              item.email_verified
                ? {
                    ...props,
                    icon: 'mdi-check-circle',
                    color: 'success'
                  }
                : {
                    ...props,
                    icon: 'mdi-clock',
                    color: 'warning'
                  }
            "
          />
        </template>
        {{ item.email_verified ? 'Verified' : 'Not verified' }}
      </v-tooltip>
    </template>
    <template #text>
      <div v-show="active">
        <div v-if="item.organisation">Organisation: {{ item.organisation }}</div>
        <div class="font-italic">
          {{ item.motive ?? 'No motive provided.' }}
        </div>
      </div>
    </template>
    <v-divider></v-divider>
    <template #actions>
      <ItemDateChip :date="item.created_on" icon="mdi-calendar-clock" />
      <v-btn
        class="ms-auto"
        color="primary"
        text="Invite"
        prepend-icon="mdi-account-plus"
        variant="plain"
        @click.stop="emit('invite', item)"
      />
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { PendingUserRequest } from '@/api'
import ItemDateChip from '@/components/toolkit/ItemDateChip.vue'
import { useToggle } from '@vueuse/core'

defineProps<{ item: PendingUserRequest }>()
const selected = defineModel<boolean>('selected')

const [active, toggleActive] = useToggle(false)

const emit = defineEmits<{ invite: [item: PendingUserRequest] }>()
</script>

<style scoped></style>
