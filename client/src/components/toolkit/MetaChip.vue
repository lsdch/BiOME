<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip
        label
        color="#777"
        variant="tonal"
        v-bind="{ ...props, ...$attrs }"
        :text="Meta.toString(meta)"
      >
        <template v-slot:prepend>
          <v-icon color="#777" class="mr-3" size="small" :icon="Meta.icon(meta)" />
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
              ? DateTime.fromJSDate(meta.modified).toLocaleString(DateTime.DATE_FULL, {
                  locale: 'en-gb'
                })
              : undefined
          "
          :subtitle="
            meta.modified
              ? `Last updated at ${DateTime.fromJSDate(meta.created).toLocaleString(
                  DateTime.TIME_24_SIMPLE,
                  {
                    locale: 'en-gb'
                  }
                )}`
              : 'Never updated'
          "
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

type Props = {
  meta: Meta
  iconColor?: string
}
withDefaults(defineProps<Props>(), {
  iconColor: '#777'
})
</script>

<style scoped lang="scss"></style>
