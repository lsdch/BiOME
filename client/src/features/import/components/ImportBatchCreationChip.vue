<template>
  <v-menu>
    <template #activator="{ props }">
      <v-chip
        label
        v-bind="props"
        class="ma-1"
        variant="tonal"
        size="small"
        :prepend-icon="batch.completed_at ? 'mdi-content-save-check' : 'mdi-content-save-move'"
      >
        {{
          batch.completed_at
            ? DateTime.fromJSDate(batch.completed_at).toLocaleString(DateTime.DATETIME_SHORT)
            : DateTime.fromJSDate(batch.created_at).toLocaleString(DateTime.DATETIME_SHORT)
        }}
      </v-chip>
    </template>
    <v-card>
      <v-list>
        <v-list-item
          :subtitle="
            DateTime.fromJSDate(batch.created_at).toRelative({ locale: 'en' }) ?? undefined
          "
        >
          <template #title>
            Submitted by <UserChip :user="batch.created_by_user" size="small" />
          </template>
        </v-list-item>
        <v-list-item
          v-if="batch.completed_at"
          :subtitle="
            DateTime.fromJSDate(batch.completed_at).toRelative({ locale: 'en' }) ?? undefined
          "
        >
          <template #title>
            Materialized by <UserChip :user="batch.completed_by_user" size="small" />
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { ImportBatchListItem } from '@/api'
import UserChip from '@/features/users/components/UserChip'
import { DateTime } from 'luxon'

defineProps<{
  batch: ImportBatchListItem
}>()
</script>

<style scoped lang="scss"></style>
