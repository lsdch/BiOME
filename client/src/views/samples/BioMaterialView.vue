<template>
  <CRUDTable
    class="fill-height"
    entity-name="Bio material"
    :headers
    :toolbar="{ title: 'Bio material', icon: 'mdi-package-variant' }"
    :fetch-items="SamplesService.listBioMaterial"
    :delete="({ code }: BioMaterial) => SamplesService.deleteBioMaterial({ path: { code } })"
    append-actions
  >
    <template #item.type="{ value }: { value: BioMaterialType }">
      <v-icon
        v-bind="
          {
            Internal: { icon: 'mdi-cube-scan', color: 'primary' },
            External: { icon: 'mdi-open-in-new', color: 'warning' }
          }[value]
        "
        :title="value"
      />
    </template>
    <template #item.identification.taxon="{ value: taxon }: { value: Taxon }">
      <TaxonChip :taxon />
    </template>
    <template #item.identification.identified_by="{ value }: { value: PersonInner }">
      <v-chip :text="value.full_name" />
    </template>
    <template #item.identification.identified_on="{ value }">
      {{ DateWithPrecision.format(value) }}
    </template>
    <template #expanded-row-inject="{ item }">
      <v-row v-if="item.external" class="ma-0">
        <v-col v-if="item.external.original_taxon != undefined">
          <v-card>
            <v-list>
              <v-list-subheader>Original identification</v-list-subheader>
              <v-chip :text="item.external.original_taxon"></v-chip>
            </v-list>
          </v-card>
        </v-col>
        <v-col>
          <v-card>
            <v-list density="compact">
              <v-list-subheader>Collection </v-list-subheader>
              <v-list-item :title="item.external.archive.collection">
                <template #append>
                  <v-btn
                    icon="mdi-link-variant"
                    :href="item.external.original_link"
                    size="small"
                    variant="tonal"
                  />
                </template>
                <v-chip
                  v-for="voucher in item.external.archive.vouchers"
                  size="small"
                  :text="voucher"
                  class="mx-1"
                />
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
        <v-col>
          <v-card>
            <v-list density="compact">
              <v-list-subheader>Content</v-list-subheader>
              <v-list-item lines="one" :subtitle="item.external.content_description">
                <template #title>
                  Quantity: <v-chip :text="item.external.quantity" size="small" />
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { BioMaterial, BioMaterialType, PersonInner, SamplesService, Taxon } from '@/api'
import { DateWithPrecision } from '@/api/adapters'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { DateTime } from 'luxon'

const headers: CRUDTableHeader[] = [
  { key: 'code', title: 'Code' },
  { key: 'type', title: 'Type', width: 0, align: 'center' },
  {
    key: 'identification',
    title: 'Identification',
    align: 'center',
    children: [
      { key: 'identification.taxon', title: 'Taxon', align: 'center' },
      { key: 'identification.identified_by', title: 'Done by', align: 'center' },
      {
        key: 'identification.identified_on',
        title: 'Date',
        align: 'end',
        sort(a: DateWithPrecision, b: DateWithPrecision) {
          return (DateWithPrecision.toDateTime(b) ?? 0) > (DateWithPrecision.toDateTime(a) ?? 0)
            ? 1
            : -1
        }
      }
    ]
  }
]
</script>

<style scoped lang="scss"></style>
