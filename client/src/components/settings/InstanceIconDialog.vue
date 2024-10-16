<template>
  <v-dialog v-model="model" :fullscreen="smAndDown" :max-width="1000">
    <v-card flat :rounded="0">
      <v-toolbar flat dark>
        <v-card-title>
          <v-icon>mdi-image</v-icon>
          App icon
        </v-card-title>
        <v-spacer />
        <v-btn color="secondary" text="Cancel" @click="model = false" />
        <v-btn
          color="primary"
          prepend-icon="mdi-floppy"
          text="Save"
          variant="flat"
          @click="saveIcon"
        />
      </v-toolbar>
      <v-container class="fill-height pa-0" fluid>
        <v-row class="fill-height">
          <v-col class="d-flex align-center" cols="12" sm="8">
            <div class="d-flex flex-column align-center mx-auto" style="max-width: 300px">
              <v-responsive :aspect-ratio="1" :max-width="300" :max-height="300">
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
          <v-divider vertical></v-divider>
          <v-col>
            <div class="h-100 d-flex flex-column align-center justify-space-between">
              <InstanceIconPreviews
                :result="result"
                :direction="smAndUp ? 'vertical' : 'horizontal'"
              />
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { SettingsService } from '@/api'
import { errorFeedback } from '@/api/responses'
import { useObjectUrl } from '@vueuse/core'
import mime from 'mime'
import { computed, ref } from 'vue'
import { CircleStencil, Coordinates, Cropper, CropperResult } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import { useDisplay } from 'vuetify'
import { useInstanceSettings } from '.'
import InstanceIconPreviews from './InstanceIconPreviews.vue'

const { smAndDown, smAndUp } = useDisplay()
const { ICON_PATH } = useInstanceSettings()

const emit = defineEmits<{ uploaded: [] }>()

const defaultCoordinates: Coordinates = {
  width: 500,
  height: 500,
  top: 100,
  left: 100
}

const model = defineModel<boolean>({ default: false })

const result = ref<CropperResult>()
function updatePreview(c: CropperResult) {
  result.value = c
}

const imgFile = ref<File | undefined>(undefined)
const imgTmpURL = useObjectUrl(imgFile)
const imgSrc = computed(() => {
  return imgTmpURL.value ?? ICON_PATH
})

// Icon mime type
const mimeType = computed(
  () => imgFile.value?.type ?? (mime.getType(ICON_PATH) || 'application/octet-stream')
)

function saveIcon() {
  result.value?.canvas?.toBlob(async (blob) => {
    if (blob !== null) {
      const file = new File([blob], 'icon', { type: mimeType.value })
      await SettingsService.setAppIcon({ body: { icon: file } })
        .then(errorFeedback('Failed to upload new icon'))
        .then(() => emit('uploaded'))
        .finally(() => (model.value = false))
    }
  }, mimeType.value)
}
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
