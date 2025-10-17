<template>
  <CardDialog title="Export occurrence data">
    <template #activator="props">
      <slot name="activator" v-bind="props"></slot>
    </template>
    <v-divider />
    <v-card-text>
      <v-select
        v-model="targetFeeds"
        label="Include feeds"
        :items="feeds"
        item-title="name"
        item-value="id"
        multiple
        chips
        closable-chips
        class="mt-2 required"
      />
      <v-select
        v-model="format"
        label="Export as"
        :items="['CSV', 'JSON']"
        hide-details
        :prepend-inner-icon="format == 'CSV' ? 'mdi-table' : 'mdi-code-braces'"
        :rounded="format == 'CSV' ? 'b-0' : undefined"
      >
        <template #selection="{ item }">
          <code> {{ item.value }}</code>
        </template>
        <template #item="{ item, props }">
          <v-list-item
            v-bind="props"
            :prepend-icon="item.value == 'CSV' ? 'mdi-table' : 'mdi-code-braces'"
          >
            <template #title>
              <code> {{ item.value }}</code>
            </template>
          </v-list-item>
        </template>
      </v-select>
      <v-card v-if="format == 'CSV'" class="bg-main mb-5 rounded-t-0">
        <v-card-text>
          <CSVDelimiterPicker v-model="optionsCSV.delimiter" />
          <div class="d-flex ga-5">
            <v-checkbox label="Force quotes" v-model="optionsCSV.quotes" />
            <CSVQuotePicker v-model="optionsCSV.quoteChar" />
          </div>
        </v-card-text>
      </v-card>

      <v-divider class="my-5" />

      <v-text-field v-model="filename" label="Filename" class="required" :suffix>
        <template #append>
          <v-btn
            prepend-icon="mdi-download"
            :text="`Export`"
            :loading="anyLoading || exporting"
            @click="executeExport()"
          />
        </template>
      </v-text-field>
    </v-card-text>
  </CardDialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import { useDataFeeds } from '@/features/cartography/components/data-feeds'
import CSVDelimiterPicker from '@/components/toolkit/ui/exports/CSVDelimiterPicker.vue'
import { DateTime } from 'luxon'
import { useToggle, useWebWorkerFn } from '@vueuse/core'
import CSVQuotePicker from '@/components/toolkit/ui/exports/CSVQuotePicker.vue'
import { useExportOptions } from '@/composables/data_exports'

const { addDataFeed, feeds, data, resetAll, context, contextEnabled, applyContext, anyLoading } =
  useDataFeeds()

const { options: optionsCSV } = useExportOptions()

type ExportFormat = 'CSV' | 'JSON'
const format = ref<ExportFormat>('CSV')

const targetFeeds = ref(feeds.value.map(({ id }) => id))

const filename = ref(`occurrences_data_${DateTime.now().toFormat('yyyy-MM-dd')}`)

const suffix = computed(() => {
  return format.value === 'JSON' ? '.json' : '.csv'
})

const [exporting, toggleExporting] = useToggle(false)

async function executeExport() {
  toggleExporting(true)
  const useFeeds = feeds.value.filter((feed) => targetFeeds.value.includes(feed.id))
  if (format.value === 'CSV') {
    // exportCSV()
  } else {
    const file = `${filename.value}${suffix.value}`
    const exportData = useFeeds.map(({ id, name }) => ({
      feed: name,
      data: data.get(id)?.data.value
    }))
    const jsonString = JSON.stringify(exportData, null, 2)
    const url = await blobify(jsonString)
    downloadObjectURL(url, file)
    toggleExporting(false)
  }
}

const { workerFn: blobify } = useWebWorkerFn((jsonString: any) => {
  const blob = new Blob([jsonString], { type: 'application/json' })
  return URL.createObjectURL(blob)
})

function downloadObjectURL(url: string, fileName: string) {
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
</script>

<style scoped lang="scss"></style>
