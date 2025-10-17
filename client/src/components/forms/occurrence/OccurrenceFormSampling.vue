<template>
  <v-card
    title="Sampling"
    :subtitle="site ? 'At site' : 'Waiting for site definition'"
    class="small-card-title"
    prepend-icon="mdi-package-down"
    flat
  >
    <template #subtitle>
      <v-chip
        v-if="hasID(sampling)"
        label
        text="From database"
        class="font-monospace"
        color="purple"
        prepend-icon="mdi-link"
        size="x-small"
        variant="flat"
      />
      <v-chip
        v-else-if="!!sampling"
        label
        text="New sampling"
        class="font-monospace"
        color="success"
        prepend-icon="mdi-plus"
        size="x-small"
        variant="flat"
      />

      <v-card-subtitle v-else>
        {{
          hasID(site)
            ? 'Pick or register sampling'
            : !!site
              ? 'Register sampling at new site'
              : 'Waiting for site definition'
        }}
      </v-card-subtitle>
    </template>
    <template #append v-if="site">
      <v-btn
        v-show="!!sampling && !hasID(sampling) && !showEdit"
        icon="mdi-pencil"
        variant="tonal"
        size="small"
        @click="toggleEdit(true)"
      />
      <SamplingFormDialog
        v-show="!sampling || hasID(sampling) || showEdit"
        v-model:dialog="dialog"
        :site
        btn-text="Save"
        subtitle="Saving does not immediately persist the sampling in the DB"
        @submit="updateSampling"
      >
        <template #activator="{ props }">
          <v-btn
            v-bind="{
              ...props,
              ...(!sampling || hasID(sampling)
                ? {
                    text: 'New sampling',
                    prependIcon: 'mdi-plus'
                  }
                : {
                    text: 'Edit new sampling',
                    prependIcon: 'mdi-pencil'
                  })
            }"
            variant="tonal"
            rounded="md"
          />
        </template>
      </SamplingFormDialog>
    </template>
    <v-card-text>
      <SiteSamplingPicker
        v-if="hasID(site)"
        :siteCode="site.code"
        @update:model-value="updateSampling"
      />
      <v-list v-if="!!sampling">
        <SamplingListItems :sampling="sampling" />
      </v-list>
      <!-- <SamplingSelectCarousel
        v-else-if="hasID(site) && site.samplings"
        :samplings="site.samplings"
        @update:model-value="updateSampling"
      /> -->
    </v-card-text>
    <template #actions v-if="showEdit && !!sampling">
      <v-spacer />
      <v-btn text="Cancel" @click="toggleEdit(false)" />
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { Sampling, Site } from '@/api'
import SamplingListItems from '@/features/occurrences/components/sampling/SamplingListItems.vue'
import SamplingFormDialog from '@/components/forms/SamplingFormDialog.vue'
import { hasID } from '@/lib/db'
import { SamplingModel } from '@/models'
import { EventModel } from '@/models/event'
import { SamplingFormModel } from '@/models/sampling'
import { useToggle } from '@vueuse/core'
import { watch } from 'vue'
import SamplingSelectCarousel from './SamplingSelectCarousel.vue'
import { SiteFormModel } from '@/models/site'
import SiteSamplingPicker from './SiteSamplingPicker.vue'

const sampling = defineModel<Sampling | SamplingFormModel>()
const dialog = defineModel<boolean>('dialog')
const props = defineProps<{ site?: Site | SiteFormModel }>()

const [showEdit, toggleEdit] = useToggle(false)

watch(
  () => props.site,
  () => updateSampling(undefined)
)

function updateSampling(s: Sampling | SamplingModel.SamplingFormModel | undefined) {
  sampling.value = s
  dialog.value = false
  toggleEdit(!s)
}
</script>

<style scoped lang="scss"></style>
