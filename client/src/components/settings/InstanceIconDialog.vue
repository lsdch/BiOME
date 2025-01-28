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
          :loading="isPending"
        />
      </v-toolbar>
      <v-container height="100%" class="overflow-y-scroll" fluid>
        <v-row v-if="error" class="mb-0">
          <v-col cols="12">
            <v-alert color="error" icon="mdi-alert"> Failed to upload new icon </v-alert>
          </v-col>
        </v-row>
        <v-row class="h-100">
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
                accept="image/png,image/jpeg"
                color="primary"
                label="Upload"
                prepend-icon=""
                prepend-inner-icon="mdi-upload"
                variant="filled"
                show-size
                hide-details
                flat
                :rounded="0"
              />
            </div>
          </v-col>
          <v-divider vertical></v-divider>
          <v-col cols="12" sm="4">
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
import { setAppIconMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'
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

const { feedback } = useFeedback()

const { mutateAsync, error, isPending } = useMutation({
  ...setAppIconMutation(),
  onError: (error) => {
    feedback({ message: 'Failed to upload new icon', type: 'error' })
    console.error(error)
  },
  onSuccess: () => {
    feedback({ message: 'Icon update', type: 'success' })
    emit('uploaded')
    model.value = false
  }
})

function saveIcon() {
  result.value?.canvas?.toBlob(async (blob) => {
    if (blob !== null) {
      const file = new File([blob], 'icon', { type: mimeType.value })
      await mutateAsync({ body: { icon: file } })
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
