<template>
  <v-bottom-sheet v-model="open" :inset="mdAndUp" content-class="rounded-0">
    <v-card :rounded="false">
      <template #prepend>
        <LinkIconGBIF v-if="taxon.GBIF_ID" :GBIF_ID="taxon.GBIF_ID" variant="text" />
        <FTaxonStatusIndicator v-else :status="taxon.status" />
      </template>

      <template #append>
        <v-btn variant="text" icon="mdi-close" @click="open = false" />
      </template>

      <template #title>
        {{ taxon.name }}
      </template>
      <template #subtitle>
        {{ taxon.authorship }}
      </template>

      <template #text>
        <div class="d-flex">
          <div class="flex-grow-1">
            <v-hover>
              <template #default="{ isHovering, props }">
                <div class="w-auto" v-bind="props">
                  <v-text-field
                    class="shrink"
                    v-model="taxon.code"
                    label="Code"
                    variant="plain"
                    readonly
                  >
                    <template #append>
                      <v-btn
                        color="primary"
                        v-if="isHovering"
                        icon="mdi-pencil-circle-outline"
                        density="compact"
                        variant="plain"
                      />
                    </template>
                  </v-text-field>
                </div>
              </template>
            </v-hover>
          </div>
          <v-spacer></v-spacer>
          <div class="flex-grow-0">
            <v-chip :text="taxon.rank" variant="outlined" class="ma-3" />
            <v-chip
              :text="taxon.status"
              variant="outlined"
              :color="taxonStatusIndicatorProps(taxon.status).color"
            />
          </div>
        </div>

        <div v-if="taxon.anchor" class="text-body-2">
          <v-icon icon="mdi-pin" color="warning" />
          This taxon was the origin of an import from GBIF on
          {{ moment(taxon.meta.created).format('DD/MM/YYYY') }}
        </div>

        <div>{{ taxon.comment }}</div>

        <v-divider class="my-3" />

        <v-list-subheader> Lineage </v-list-subheader>
        <div class="lineage" v-if="relatives">
          <v-skeleton-loader type="chip@5">
            <template
              v-for="(v, i) in Object.values(relatives.lineage).filter((v) => Boolean(v))"
              :key="i"
            >
              <v-btn
                color="primary"
                class="text-body-2"
                variant="text"
                :text="v?.name"
                @click="emit('navigate', v!)"
              />
              <v-icon>mdi-chevron-right</v-icon>
            </template>
            <span class="text-body-2 px-4">
              {{ taxon.name }}
            </span>
          </v-skeleton-loader>
        </div>

        <v-list-subheader>
          Descendants
          <v-chip color="primary" :text="`${taxon.children_count}`" :rounded="100" size="small" />
        </v-list-subheader>
        <div class="descendants">
          <v-skeleton-loader type="chip@5">
            <v-chip
              v-for="c in relatives?.children"
              :key="c.id"
              class="ma-2"
              @click="emit('navigate', c)"
            >
              {{ c.name }}
            </v-chip>
          </v-skeleton-loader>
        </div>
      </template>

      <v-divider />

      <template #actions>
        <div>
          <ItemDateChip v-if="taxon.meta?.created" icon="created" :date="taxon.meta.created" />
          <ItemDateChip v-if="taxon.meta?.modified" icon="updated" :date="taxon.meta.modified" />
        </div>
        <v-spacer />
        <div v-if="isGranted('Admin')">
          <v-btn
            color="error"
            text="Delete"
            prepend-icon="mdi-delete-outline"
            @click="deleteTaxon(taxon)"
          />
          <v-btn
            v-if="extensibleRanks.includes(taxon.rank)"
            color="primary"
            text="Add descendant"
            prepend-icon="mdi-arrow-decision"
            @click="emit('add-child', taxon)"
          />
        </div>
      </template>
    </v-card>
  </v-bottom-sheet>
</template>

<script setup lang="ts">
import { Taxon, TaxonomyService, TaxonRank, TaxonWithLineage, TaxonWithRelatives } from '@/api'
import { handleErrors } from '@/api/responses'
import { useAppConfirmDialog } from '@/composables'
import { useFeedback } from '@/stores/feedback'
import { ref, watch } from 'vue'
import { useDisplay } from 'vuetify'
import ItemDateChip from '../toolkit/ItemDateChip.vue'
import LinkIconGBIF from './LinkIconGBIF.vue'
import { useUserStore } from '@/stores/user'
import moment from 'moment'
import { FTaxonStatusIndicator, taxonStatusIndicatorProps } from './functionals'

const { mdAndUp } = useDisplay()
const { isGranted } = useUserStore()

const extensibleRanks: TaxonRank[] = ['Order', 'Family', 'Genus', 'Species']

const taxon = defineModel<Taxon>({ required: true })
const open = defineModel<boolean>('open')

const relatives = ref<TaxonWithLineage>()
const loading = ref(false)

const emit = defineEmits<{
  'add-child': [parent: Taxon]
  deleted: [taxon: TaxonWithRelatives]
  navigate: [target: Taxon]
}>()

watch(
  taxon,
  async (taxon) => {
    relatives.value = undefined
    relatives.value = await fetch(taxon)
  },
  { immediate: true }
)

async function fetch(taxon: Taxon) {
  loading.value = true
  const data = await TaxonomyService.getTaxon({ path: { code: taxon.code } }).then(
    handleErrors((err) => console.error('Failed to fetch taxon', err))
  )
  loading.value = false
  return data
}

const { askConfirm } = useAppConfirmDialog()
const { feedback } = useFeedback()

async function deleteTaxon(taxon: Taxon) {
  askConfirm({
    title: `Delete taxon ${taxon.name}?`,
    message: 'All descendants will also be deleted'
  }).then(async ({ isCanceled }) => {
    if (isCanceled) return
    const { data, error } = await TaxonomyService.deleteTaxon({ path: { code: taxon.code } })
    if (error !== undefined) {
      feedback({ type: 'error', message: 'Failed to delete taxon' })
      return
    }
    emit('deleted', data)
    open.value = false
    feedback({ type: 'success', message: `Taxon ${taxon.name} was successfully deleted` })
  })
}
</script>

<style scoped>
.descendants {
  max-height: 50dvh;
  overflow-y: scroll;
}
</style>
