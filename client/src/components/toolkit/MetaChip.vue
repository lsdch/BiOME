<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip
        label
        color="#777"
        variant="tonal"
        v-bind="{ ...props, ...$attrs }"
        :text="DateTime.fromJSDate(meta.last_updated).toRelative({ locale: 'en-gb' })?.toString()"
      >
        <template v-slot:prepend>
          <v-icon color="#777" class="mr-3" size="small" :icon="icon" />
        </template>
      </v-chip>
    </template>
    <v-card :min-width="350" class="bg-surface-light" density="compact">
      <v-list class="bg-surface-light" density="compact">
        <v-list-item
          :title="
            DateTime.fromJSDate(meta.created).toLocaleString(DateTime.DATE_FULL, {
              locale: 'en-gb'
            })
          "
          :subtitle="`Created at ${DateTime.fromJSDate(meta.created).toLocaleString(
            DateTime.TIME_24_SIMPLE,
            {
              locale: 'en-gb'
            }
          )}`"
          prepend-icon="mdi-content-save"
          slim
        >
          <template #append>
            <v-chip v-if="meta.created_by" :text="meta.created_by.name" size="small" />
            <v-chip text="No author" size="small" />
          </template>
        </v-list-item>
        <v-divider />
        <v-list-item
          :title="
            meta.modified
              ? DateTime.fromJSDate(meta.modified).toLocaleString({}, { locale: 'en-gb' })
              : undefined
          "
          :subtitle="meta.modified ? 'Last updated' : 'Never updated'"
          prepend-icon="mdi-update"
          slim
        />
      </v-list>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Meta } from '@/api'
import { DateTime } from 'luxon'
import { computed } from 'vue'

type UpdateIcon = 'mdi-update'
type CreatedIcon = 'mdi-content-save'
type Icon = UpdateIcon | CreatedIcon | string

type Props = {
  meta: Meta
  iconColor?: string
}
const props = withDefaults(defineProps<Props>(), {
  iconColor: '#777'
})

const icon = computed<string>(() => {
  return props.meta.modified ? 'mdi-update' : 'mdi-content-save'
})
</script>

<style scoped lang="scss"></style>
