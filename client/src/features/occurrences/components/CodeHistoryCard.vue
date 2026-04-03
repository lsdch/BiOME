<template>
  <v-card class="rounded-t-0">
    <v-list>
      <v-list-subheader>Code history</v-list-subheader>
      <v-list-item v-for="item in sortedHistory" :title="item.code" @click="onClick">
        <template #title="{ title }">
          <span class="font-monospace">{{ title }}</span>
        </template>
        <template #subtitle>
          <span class="text-caption text-muted">
            {{
              DateTime.fromJSDate(item.time).toLocaleString(DateTime.DATETIME_MED, {
                locale: 'en-gb'
              })
            }}
          </span>
        </template>
        <v-overlay :model-value="copied === item.code" contained content-class="w-100 h-100">
          <div class="bg-success w-100 h-100 d-flex align-center px-4 opacity-80">
            <v-icon>mdi-content-copy</v-icon>
            Copied to clipboard
          </div>
        </v-overlay>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { CodeHistory } from '@/api'
import { useClipboard, useThrottleFn } from '@vueuse/core'
import { DateTime } from 'luxon'
import { computed, ref } from 'vue'

const props = defineProps<{
  codeHistory: CodeHistory[]
}>()

const sortedHistory = computed(() => {
  return props.codeHistory.toSorted((a, b) => b.time.getTime() - a.time.getTime())
})

const { copy } = useClipboard()
const copied = ref<string>()

const onClick = useThrottleFn(async (ev: Event) => {
  const target = ev.target as HTMLElement
  await copy(target.textContent || '')
  copied.value = target.textContent || undefined
  setTimeout(() => {
    copied.value = undefined
  }, 1000)
}, 1000)
</script>

<style scoped lang="scss"></style>
