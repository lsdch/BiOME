<template>
  <v-img
    id="instance-icon"
    :src="ICON_PATH"
    :width="120"
    :height="120"
    alt="alt"
    rounded="circle"
    :class="['border-lg', iconHover ? 'border-primary border-opacity-100' : 'border-opacity-10']"
  >
    <v-overlay
      v-model="iconHover"
      @click="emit('click')"
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

  <InstanceIconDialog v-model="dialogOpen" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useInstanceSettings } from '.'
import InstanceIconDialog from './InstanceIconDialog.vue'

const { ICON_PATH } = useInstanceSettings()

const iconHover = ref(false)
const dialogOpen = ref(false)

const emit = defineEmits<{
  click: []
}>()
</script>

<style scoped>
.avatar-hover {
  border: 3px solid rgb(0, 112, 177);
}
</style>
