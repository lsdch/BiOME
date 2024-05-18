<template>
  <v-dialog v-model="model.open" :fullscreen="smAndDown" :max-width="1000">
    <v-card flat :rounded="0">
      <v-toolbar flat dark>
        <v-card-title>
          <v-icon>mdi-image</v-icon>
          App icon
        </v-card-title>
        <v-spacer />
        <v-btn color="secondary" text="Cancel" @click="model.open = false" />
        <v-btn
          color="primary"
          prepend-icon="mdi-floppy"
          text="Save"
          variant="flat"
          @click="saveIcon"
        />
      </v-toolbar>
      <v-container fluid>
        <v-row>
          <v-col cols="12" md="8">
            <div class="d-flex flex-column mx-auto" style="max-width: 500px">
              <v-responsive :aspect-ratio="1" :max-width="500" :max-height="500">
                <cropper
                  class="cropper icon-cropper"
                  @change="updatePreview"
                  auto-zoom
                  :src="imgSrc"
                  :default-visible-area="defaultCoordinates"
                  :stencil-component="CircleStencil"
                />
              </v-responsive>
              <v-file-input
                v-model="imgFile"
                class="w-100 cursor-pointer"
                accept="image/*"
                color="primary"
                label="Upload new image"
                prepend-icon=""
                prepend-inner-icon="mdi-upload"
                variant="solo-filled"
                show-size
                single-line
                hide-details
                flat
                :rounded="0"
              />
            </div>
          </v-col>
          <v-col>
            <div class="h-100 d-flex flex-column align-center justify-space-between">
              <InstanceIconPreviews
                :result="result"
                :direction="mdAndUp ? 'vertical' : 'horizontal'"
              />
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CircleStencil, Coordinates, Cropper, CropperResult } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import { useDisplay } from 'vuetify'
import InstanceIconPreviews from './InstanceIconPreviews.vue'
import { computed, watch } from 'vue'
import { SettingsService } from '@/api'
import mime from 'mime'

const { smAndDown, mdAndUp } = useDisplay()

const defaultCoordinates: Coordinates = {
  width: 500,
  height: 500,
  top: 100,
  left: 100
}

const model = defineModel<{
  open: boolean
  iconSrc?: string
}>({ default: { open: false } })

const result = ref<CropperResult>()
function updatePreview(c: CropperResult) {
  result.value = c
}

const imgFile = ref<File | undefined>(undefined)
const uploadedImg = ref<string | undefined>(undefined)

watch(imgFile, () => {
  if (uploadedImg.value) {
    URL.revokeObjectURL(uploadedImg.value)
  }
  if (imgFile.value) {
    uploadedImg.value = URL.createObjectURL(imgFile.value)
  } else {
    uploadedImg.value = undefined
  }
})

const imgSrc = computed(() => {
  return uploadedImg.value ?? model.value.iconSrc
})

function saveIcon() {
  result.value?.canvas?.toBlob((blob) => {
    if (blob !== null) {
      const file = new File([blob], 'icon', { type: mimeType.value })
      SettingsService.setAppIcon({ formData: { icon: file } })
    }
  }, mimeType.value)
}

const mimeType = computed(
  () =>
    imgFile.value?.type ?? (mime.getType(model.value.iconSrc ?? '') || 'application/octet-stream')
)
</script>

<style lang="scss">
@use 'vuetify';

.icon-cropper {
  width: 500px;
  aspect-ratio: 1 / 1;
  .vue-advanced-cropper__background {
    background: #252525;
  }
}

.crop-container {
  max-width: 500px;
  max-height: 500px;
  display: flex;
  justify-content: center;
  overflow: hidden;
}

.v-file-input {
  * {
    cursor: pointer;
  }
  i {
    color: white;
  }
  .v-field {
    background: rgb(var(--v-theme-primary)) !important;
  }
  label,
  .v-field__input {
    color: white !important;
  }
}
</style>
