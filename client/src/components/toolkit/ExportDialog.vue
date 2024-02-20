<template>
  <v-dialog v-model="dialog" max-width="500px">
    <v-card>
      <v-toolbar dark dense flat>
        <v-card-title> Export {{ items.length }} items as CSV/TSV </v-card-title>
        <template v-slot:append>
          <v-btn color="grey" @click="dialog = false" icon="mdi-close"></v-btn>
        </template>
      </v-toolbar>
      <v-card-text>
        <v-form @submit.prevent>
          <v-text-field v-model="filename" label="Filename" :suffix="suffix" />
          <v-row>
            <v-col cols="12" sm="">
              <v-checkbox label="Quote items" v-model="options.quotes" color="primary"></v-checkbox>
            </v-col>
            <v-col cols="12" sm="">
              <v-select
                label="Quote char"
                :items="quoteChars"
                v-model="options.quoteChar"
                item-value="value"
                :item-props="(item) => item"
                :disabled="!options.quotes"
              />
            </v-col>
          </v-row>
          <v-select
            label="Delimiter"
            :items="delimiters"
            v-model="options.delimiter"
            item-value="value"
            :item-props="(item) => item"
          />
        </v-form>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions class="d-flex justify-center">
        <v-btn
          variant="text"
          color="primary"
          v-bind="button"
          @click="revokeURL"
          prepend-icon="mdi-download"
          text="Download"
          :loading="loading"
        />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts" generic="ItemType extends {}">
import { flatten } from 'flat'
import moment from 'moment'
import CSVEngine from 'papaparse'
import { computed, ref, watch } from 'vue'

const dialog = defineModel<boolean>()
const props = defineProps<{ items: ItemType[]; namePrefix: string }>()

function defaultOptions() {
  return {
    delimiter: '\t',
    quotes: true,
    quoteChar: '"'
  }
}

const options = ref(defaultOptions())
const filename = ref(generateFilename())

const delimiters = [
  { title: '\\t', subtitle: 'Tab', value: '\t' },
  { title: ',', subtitle: 'Comma', value: ',' },
  { title: ';', subtitle: 'Semicolon', value: ';' }
]

const quoteChars = [
  { title: '"', subtitle: 'Double', value: '"' },
  { title: "'", subtitle: 'Single', value: "'" }
]

function revokeURL() {
  URL.revokeObjectURL(button?.value?.href)
}

watch(props, () => {
  revokeURL()
  options.value = defaultOptions()
})

const suffix = computed(() => (options.value.delimiter === '\t' ? '.tsv' : '.csv'))

function generateFilename() {
  return `${props.namePrefix}_${moment().format('Y-MM-DD')}`
}

const emit = defineEmits<{
  ready: []
}>()

const csvString = ref('')
const loading = ref(true)

watch(() => props.items, unparse, { immediate: true })
watch(() => options.value, unparse, { immediate: true, deep: true })

async function unparse() {
  loading.value = true
  csvString.value = await new Promise((resolve) => {
    const res = CSVEngine.unparse(
      props.items.map((item) => flatten(item)),
      {
        ...options.value,
        quotes(value) {
          return options.value.quotes && !['boolean', 'number', 'bigint'].includes(typeof value)
        }
      }
    )
    resolve(res)
  })
  const blob = new Blob([csvString.value], { type: 'text/csv;charset=utf8' })
  button.value = {
    href: URL.createObjectURL(blob),
    download: filename.value + suffix.value
  }
  emit('ready')
  loading.value = false
}

const button = ref({ href: '', download: '' })
</script>

<style scoped></style>
