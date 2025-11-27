<template>
  <CardDialog :title v-bind="props">
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData ?? {}" />
    </template>
    <div v-if="(occurrences?.length ?? 0) > 10" class="d-flex align-center ga-2">
      <v-text-field
        v-model="search.term"
        class="mx-5"
        hide-details
        label="Search"
        clearable
        density="compact"
      />
    </div>
    <CRUDTable :items="occurrences" entity-name="Occurrences" :headers :search>
      <template #item.code="{ item, value }: { item: OccurrenceAtSite; value: string }">
        <div class="d-flex justify-space-between align-center">
          <RouterLink
            :to="{
              name: 'occurrence-item',
              params: { code: item.code }
            }"
            target="_blank"
          >
            <span class="text-wrap">{{ CodeIdentifier.textWrap(value) }}</span>
          </RouterLink>
        </div>
      </template>
      <template #item.site.code="{ value }: { value: string }">
        <RouterLink
          :to="{
            name: 'site-item',
            params: { code: value }
          }"
          target="_blank"
        >
          <span class="text-wrap font-monospace">
            {{ CodeIdentifier.textWrap(value) }}
          </span>
        </RouterLink>
      </template>
      <template #item.site.coordinates.latitude="{ value }">
        <span class="font-monospace">
          {{ value }}
        </span>
      </template>
      <template #item.site.coordinates.longitude="{ value }">
        <span class="font-monospace">
          {{ value }}
        </span>
      </template>
      <template #item.sampling_date="{ value }">
        <span class="font-monospace">
          {{ DateWithPrecision.format(value) }}
        </span>
      </template>
      <template #item.identification="{ value: identification }">
        <IdentificationChip :identification size="small" />
      </template>
    </CRUDTable>
  </CardDialog>
</template>

<script setup lang="ts" generic="WithSite extends boolean">
import {
  CodeIdentifier,
  DateWithPrecision,
  Identification,
  OccurrenceAtSite,
  SiteItem
} from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import { computed, ref } from 'vue'
import { ComponentSlots } from 'vue-component-type-helpers'

type Occurrence = OccurrenceAtSite & { sampling_date: DateWithPrecision } & (WithSite extends true
    ? { site: SiteItem }
    : {})

const {
  title = 'Occurrences',
  occurrences,
  ...props
} = defineProps<
  {
    /**
     * If true, the site information is included with each occurrence.
     * If false, only the occurrence information is included.
     */
    withSite: WithSite
    occurrences?: Occurrence[]
  } & CardDialogProps
>()

const search = ref({
  term: undefined,
  owned: undefined
})

const headersWithSites: CRUDTableHeader<Occurrence>[] = [
  {
    title: 'Code',
    value: 'code',
    sortable: true
  },
  {
    title: 'Site',
    value: 'site.code',
    width: 0,
    sortable: true
  },
  {
    title: 'Latitude',
    value: 'site.coordinates.latitude',
    width: 0,
    sortable: true
  },
  {
    title: 'Longitude',
    value: 'site.coordinates.longitude',
    sortable: true,
    width: 0
  },
  {
    title: 'Sampl. date',
    value: 'sampling_date',
    sortable: true,
    align: 'end',
    sort: DateWithPrecision.compare
  },
  Identification.tableHeader({ key: 'identification', align: 'center' })
]

const headers = computed(() =>
  headersWithSites.filter(
    (header) =>
      props.withSite !== false ||
      typeof header.value != 'string' ||
      !header.value.startsWith('site.')
  )
)

defineSlots<ComponentSlots<typeof CardDialog>>()
</script>

<style scoped lang="scss"></style>
