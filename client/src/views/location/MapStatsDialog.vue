<template>
  <CardDialog
    v-model="dialog"
    title="Statistics"
    width="400"
    v-bind="$attrs"
    prepend-icon="mdi-poll"
  >
    <template #activator="props" v-if="$slots.activator">
      <slot name="activator" v-bind="props" />
    </template>
    <v-list>
      <v-list-item title="Sites">
        <template #append>
          <v-badge inline :content="sites?.length ?? 0" color="primary"></v-badge>
        </template>
      </v-list-item>
      <v-list-item title="Sampling events">
        <template #append>
          <v-badge
            inline
            :content="sites?.reduce((sum, s) => sum + s.samplings.length, 0) ?? 0"
            color="warning"
          ></v-badge>
        </template>
      </v-list-item>
      <v-list-item title="Occurrences">
        <template #append>
          <v-badge
            inline
            :content="
              sites
                ?.flatMap(({ samplings }) => samplings)
                .reduce((sum, s) => sum + s.occurrences.length, 0) ?? 0
            "
            color="success"
          />
        </template>
      </v-list-item>
    </v-list>
  </CardDialog>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'

const dialog = defineModel<boolean>({
  default: false
})

const { sites } = defineProps<{
  sites?: SiteWithOccurrences[]
}>()
</script>

<style scoped lang="scss"></style>
