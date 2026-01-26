<template>
  <CardDialog
    v-if="dataset?.contributors"
    :title="`Participants in dataset: ${dataset.label}`"
    prepend-icon="mdi-account-multiple"
  >
    <v-card-text>
      <div class="d-flex ga-2 flex-wrap">
        <v-chip v-for="name in sortedContributors" :key="name">{{ name }}</v-chip>
      </div>
    </v-card-text>
    <template #activator="{ props }">
      <v-chip label size="small" v-bind="props" prepend-icon="mdi-account-multiple">
        {{ dataset.contributors.length ?? 0 }} contributors
      </v-chip>
    </template>
    <template #actions>
      <span class="text-muted text-success">{{ confirmCopyMsg }}</span>
      <v-btn text="Copy" prepend-icon="mdi-content-copy" @click="makeCopy"></v-btn>
    </template>
  </CardDialog>
</template>

<script setup lang="ts">
import { OccurrenceDataset } from '@/api'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import { useClipboard } from '@vueuse/core'
import { computed, ref } from 'vue'

const { dataset } = defineProps<{ dataset: OccurrenceDataset }>()

const { copy } = useClipboard()

const confirmCopyMsg = ref()

const sortedContributors = computed(() => {
  return dataset.contributors?.toSorted() || []
})

function makeCopy() {
  copy(sortedContributors.value.join(', '))
  confirmCopyMsg.value = 'Copied to clipboard'
  setTimeout(() => {
    confirmCopyMsg.value = ''
  }, 2000)
}
</script>

<style scoped lang="scss"></style>
