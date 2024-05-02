<template>
  <v-snackbar
    v-model="isOpen"
    variant="flat"
    :color="theme.color"
    :timeout="timeout"
    timer="secondary"
  >
    <div class="text-overline mb-3 d-flex align-center">
      <v-icon class="mr-2" :icon="theme.icon" color="white" size="large" />
      {{ current?.message }}
    </div>
    <template v-slot:actions>
      <v-btn color="white" variant="text" @click="isOpen = false"> Close </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { useFeedback } from '@/stores/feedback'
import { storeToRefs } from 'pinia'
import { computed, nextTick, ref, watch } from 'vue'
import colors from 'vuetify/util/colors'

withDefaults(defineProps<{ timeout: number }>(), { timeout: 2000 })

const isOpen = ref(false)

function open() {
  nextTick(() => {
    isOpen.value = true
  })
}
watch(isOpen, (state, oldState) => {
  if (oldState && !state) {
    console.log('HAS CLOSED')
    store.next()
  }
})

const store = useFeedback()
const { current } = storeToRefs(store)
watch(
  current,
  (newCurrent) => {
    isOpen.value = false
    if (newCurrent) open()
  },
  { immediate: true }
)

const theme = computed(() => {
  switch (current.value?.type) {
    case 'error':
      return {
        color: colors.red.darken1,
        icon: 'mdi-alert-circle'
      }
    case 'info':
      return {
        color: current.value?.type,
        icon: 'mdi-information'
      }
    case 'warning':
      return {
        color: current.value?.type,
        icon: 'mdi-alert'
      }
    case 'primary':
      return {
        color: current.value?.type,
        icon: 'mdi-circle'
      }
    case 'success':
      return {
        color: current.value?.type,
        icon: 'mdi-check-circle'
      }
    default:
      return {
        color: 'secondary'
      }
  }
})
</script>

<style scoped></style>
