<template>
  <v-dialog v-model="dialog" :max-width="500">
    <v-card>
      <v-toolbar dark dense flat>
        <v-card-title> Export {{ items.length }} items </v-card-title>
        <template v-slot:append>
          <v-btn color="grey" @click="dialog = false" icon="mdi-close" />
        </template>
      </v-toolbar>
      <v-card-text>
        <v-form @submit.prevent validate-on="input" v-model="isValid">
          <v-text-field v-model="filename" label="Filename" :suffix="suffix" required />
          <v-row>
            <v-col cols="12" sm="">
              <v-checkbox label="Quote items" v-model="options.quotes" color="primary" />
            </v-col>
            <v-col cols="12" sm="">
              <CSVQuotePicker v-model="options.quoteChar" :disabled="!options.quotes" />
            </v-col>
          </v-row>
          <CSVDelimiterPicker v-model="options.delimiter" />
        </v-form>
      </v-card-text>
      <v-divider />
      <v-expansion-panels>
        <v-expansion-panel title="Optional fields" class="" :height="10">
          <template #text>
            <v-checkbox
              label="UUIDs"
              v-model="options.includeUUID"
              color="primary"
              density="compact"
            />
            <v-checkbox
              label="Modification metadata"
              v-model="options.includeModificationMetadata"
              color="primary"
              density="compact"
            />
          </template>
        </v-expansion-panel>
      </v-expansion-panels>
      <v-divider />
      <v-card-actions class="d-flex justify-center">
        <v-btn
          block
          variant="text"
          color="primary"
          :href="buttonHref"
          :download="filename + suffix"
          @click="revokeURL"
          prepend-icon="mdi-download"
          text="Download"
          :loading="loading"
          :disabled="!isValid"
        />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="ItemType extends {}">
import { flatten } from 'flat'
import CSVEngine from 'papaparse'
import { computed, ref, watch } from 'vue'
import CSVDelimiterPicker from './exports/CSVDelimiterPicker.vue'
import CSVQuotePicker from './exports/CSVQuotePicker.vue'
import { DateTime } from 'luxon'
import { useExportOptions } from '@/composables/data_exports'

const isValid = ref(null)

const dialog = defineModel<boolean>()
const props = defineProps<{ items: ItemType[]; namePrefix: string }>()
const emit = defineEmits<{ ready: [] }>()

const { options, reset } = useExportOptions()
const filename = ref(generateFilename())

function revokeURL() {
  URL.revokeObjectURL(buttonHref?.value?.href)
}

watch(props, () => {
  revokeURL()
  reset()
})

const suffix = computed(() => (options.value.delimiter === '\t' ? '.tsv' : '.csv'))

function generateFilename() {
  return `${props.namePrefix}_${DateTime.now().toFormat('yyyy-MM-dd')}`
}

const blobString = ref('')
const loading = ref(true)

watch(() => props.items, unparse, { immediate: true })
watch(() => options.value, unparse, { immediate: true, deep: true })

const buttonHref = ref()
async function unparse() {
  console.log('Unparsing data with options', options.value)
  loading.value = true

  blobString.value = await new Promise((resolve) => {
    const res = CSVEngine.unparse(
      props.items.map((item) => flatten(item)),
      {
        ...options.value,
        // exclude columns that end with '.id'
        columns: Object.keys(flatten(props.items[0])).filter((key) => {
          if ((key.endsWith('.id') || key === 'id') && !options.value.includeUUID) {
            return false
          }
          if (
            (key === 'meta' || key.startsWith('meta.')) &&
            !options.value.includeModificationMetadata
          ) {
            return false
          }
          return true
        }),
        quotes(value) {
          return options.value.quotes && !['boolean', 'number', 'bigint'].includes(typeof value)
        }
      }
    )
    resolve(res)
  })

  const blob = new Blob([blobString.value], { type: 'text/csv;charset=utf8' })
  buttonHref.value = URL.createObjectURL(blob)
  emit('ready')
  loading.value = false
}
</script>

<style scoped></style>
