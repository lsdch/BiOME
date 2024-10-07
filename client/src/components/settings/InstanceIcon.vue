<template>
  <v-img
    id="instance-icon"
    :src="`${ICON_PATH}?cacheKey=${renderKey}`"
    :width="120"
    :aspect-ratio="1"
    :max-width="120"
    alt="alt"
    rounded="circle"
    :class="['border-lg', iconHover ? 'border-primary border-opacity-100' : 'border-opacity-10']"
    :key="renderKey"
  >
    <v-overlay
      v-model="iconHover"
      @click="dialogOpen = true"
      open-on-hover
      :close-delay="200"
      class="align-center justify-center cursor-pointer font-weight-black text-white"
      scrim="primary"
      activator="parent"
      contained
    >
      Change icon
    </v-overlay>
  </v-img>

  <InstanceIconDialog v-model="dialogOpen" @uploaded="onIconUploaded" />
</template>

<script setup lang="ts">
import { ref, render } from 'vue'
import { useInstanceSettings } from '.'
import InstanceIconDialog from './InstanceIconDialog.vue'

const { ICON_PATH } = useInstanceSettings()

const iconHover = ref(false)
const dialogOpen = ref(false)

const renderKey = ref(0)
const emit = defineEmits<{ changed: [] }>()

function onIconUploaded() {
  renderKey.value++
  emit('changed')
}
</script>

<style scoped>
.avatar-hover {
  border: 3px solid rgb(0, 112, 177);
}
</style>
