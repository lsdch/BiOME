<template>
  <v-card
    ref="dropZoneRef"
    :class="['dropzone d-flex flex-column align-center justify-center', { over: isOverDropZone }]"
    :width="width"
    :height="height"
    flat
    @click="open"
  >
    <span v-if="isOverDropZone">
      <v-icon icon="mdi-file-upload" color="success" size="xx-large" />
    </span>
    <span class="d-flex flex-column align-center" v-else>
      <span class="text-overline font-weight-bold">
        <v-icon class="mr-3" icon="mdi-file-upload-outline" size="large" />
        Upload file
      </span>
      <span class="text-caption">click or drag-and-drop</span>
    </span>
    <slot name="hint">
      <v-card-actions class="dropzone-hint text-caption">{{ hint }}</v-card-actions>
    </slot>
  </v-card>
</template>

<script setup lang="ts">
import { useDropZone, useFileDialog } from '@vueuse/core'
import { ref } from 'vue'

const props = defineProps<{
  width?: string | number
  height?: number
  /**
   * List of acceptable file types
   */
  datatypes?: string[]
  multiple?: boolean
  /**
   * Help message to display at the bottom of the dropzone.
   */
  hint?: string
}>()

const emit = defineEmits<{
  upload: [files: File[]]
}>()

const { open, onChange } = useFileDialog({
  accept: props.datatypes?.join(', '),
  multiple: props.multiple
})

onChange((files) => {
  if (files == null) return
  if (files.length > 1 && !props.multiple) return
  emit('upload', Array.from(files))
})

const dropZoneRef = ref<HTMLElement>()
const { isOverDropZone } = useDropZone(dropZoneRef, {
  dataTypes: props.datatypes,
  onDrop(files) {
    if (files == null) return
    if (files.length > 1 && !props.multiple) return
    emit('upload', files)
  }
})
</script>

<style lang="scss" scoped>
@use 'vuetify';
.dropzone {
  border: 2px dashed rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-primary));
  position: relative;
  &.over {
    border-color: rgb(var(--v-theme-success));
  }
  .dropzone-hint {
    position: absolute;
    bottom: 0px;
  }
}
</style>
