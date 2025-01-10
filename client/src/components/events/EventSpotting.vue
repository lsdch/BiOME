<template>
  <v-card flat>
    <v-card-text>
      <v-list v-if="!editing">
        <v-list-item title="Taxa">
          <template #append>
            <v-btn
              variant="tonal"
              prepend-icon="mdi-pencil"
              text="Edit"
              @click="toggleEdit(true)"
            />
          </template>
          <v-chip v-for="t in spotting?.target_taxa" class="mx-1" :text="t.name" />
          <v-list-item-subtitle v-if="!spotting?.target_taxa?.length"> None </v-list-item-subtitle>
        </v-list-item>
        <v-list-item title="Comments" :subtitle="spotting?.comments || 'None'" />
      </v-list>
      <v-form v-else>
        <TaxonPicker
          label="Taxa"
          v-model="model.target_taxa"
          item-value="name"
          :ranks="TaxonRank.ranksUpTo('Family')"
          multiple
          chips
          closable-chips
          clearable
        />
        <v-textarea label="Comments" v-model="model.comments" />

        <div class="d-flex justify-end">
          <v-btn
            class="mx-1"
            color="primary"
            variant="tonal"
            text="Submit"
            :loading
            @click="submit()"
          />
          <v-btn
            class="mx-1"
            text="Cancel"
            variant="plain"
            color=""
            :disabled="loading"
            @click="reset()"
          />
        </div>
      </v-form>
      <!-- <v-btn variant="tonal" prepend-icon="mdi-pencil">Edit</v-btn> -->
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { $SpottingUpdate, Event, EventsService, Spotting, SpottingUpdate, TaxonRank } from '@/api'
import { useToggle } from '@vueuse/core'
import { ref, watch } from 'vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import { useSchema } from '../toolkit/forms/schema'

const { event } = defineProps<{ event: Event }>()
const spotting = defineModel<Spotting>('spotting', { required: true })

const [editing, toggleEdit] = useToggle(false)
const [loading, toggleLoading] = useToggle(false)

const model = ref<SpottingUpdate>(initialModel(spotting.value))

const { field, errorHandler } = useSchema($SpottingUpdate)

watch(spotting, (s) => {
  model.value = initialModel(s)
})

async function submit() {
  toggleLoading(true)
  return EventsService.updateSpotting({ path: { id: event.id }, body: model.value })
    .then(errorHandler)
    .then((updated) => {
      event.spotting = updated
      toggleEdit(false)
    })
    .finally(() => toggleLoading(false))
}

function reset() {
  model.value = initialModel(spotting.value)
  toggleEdit(false)
}

function initialModel({ target_taxa, comments }: Spotting): SpottingUpdate {
  return {
    target_taxa: target_taxa?.map(({ name }) => name) ?? null,
    comments: comments
  }
}
</script>

<style scoped lang="scss"></style>
