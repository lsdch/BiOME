<template>
  <v-bottom-sheet v-model="open" :inset="mdAndUp" content-class="rounded-0">
    <v-card :rounded="false">
      <template #prepend>
        <LinkIconGBIF v-if="taxon.GBIF_ID" :GBIF_ID="taxon.GBIF_ID" variant="text" />
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
        <div class="d-flex justify-space-between">
          <div>
            <div class="d-flex">
              Code:
              <pre class="ml-2">{{ taxon.code }}</pre>
            </div>
          </div>
          <div>
            <v-chip :text="taxon.rank" variant="outlined" class="mx-3" />
            <v-chip :text="taxon.status" variant="outlined" />
          </div>
        </div>

        <div>{{ taxon.comment }}</div>

        <v-divider class="my-3" />

        <v-list-subheader>
          Descendants
          <v-chip color="primary" :text="taxon.children_count" :rounded="100" size="small" />
        </v-list-subheader>
        <div class="descendants">
          <v-chip v-for="c in relatives?.children" :key="c.id" class="ma-2">
            {{ c.name }}
          </v-chip>
        </div>
      </template>

      <v-divider />

      <template #actions>
        <div>
          <ItemDateChip v-if="taxon.meta?.created" icon="created" :date="taxon.meta.created" />
          <ItemDateChip v-if="taxon.meta?.modified" icon="updated" :date="taxon.meta.modified" />
        </div>
        <v-spacer />
        <v-btn color="error" text="Delete" prepend-icon="mdi-delete-outline" />
        <v-btn color="primary" text="Add descendant" prepend-icon="mdi-arrow-decision" />
      </template>
    </v-card>
  </v-bottom-sheet>
</template>

<script setup lang="ts">
import { Taxon, TaxonomyService, TaxonWithRelatives } from '@/api'
import LinkIconGBIF from './LinkIconGBIF.vue'
import ItemDateChip from '../toolkit/ItemDateChip.vue'
import { handleErrors } from '@/api/responses'
import { ref, watch } from 'vue'
import { useDisplay } from 'vuetify'

const open = defineModel<boolean>()

const props = defineProps<{ taxon: Taxon }>()

const relatives = ref<TaxonWithRelatives>()

const { mdAndUp } = useDisplay()

watch(
  () => props.taxon,
  async (taxon) => {
    relatives.value = undefined
    relatives.value = await fetch(taxon)
  },
  { immediate: true }
)

async function fetch(taxon: Taxon) {
  const data = await TaxonomyService.getTaxon({ path: { code: taxon.code } }).then(
    handleErrors((err) => console.error('Failed to fetch taxon', err))
  )
  return data
}
</script>

<style scoped>
.descendants {
  max-height: 50dvh;
  overflow-y: scroll;
}
</style>
