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
          <TaxonChip v-for="taxon in model.spottings" :taxon class="mx-1" />
          <v-list-item-subtitle v-if="!model.spottings?.length"> None </v-list-item-subtitle>
        </v-list-item>
        <v-list-item title="Comments" :subtitle="model.comments || 'None'" />
      </v-list>
      <v-form v-else>
        <v-confirm-edit v-model="model">
          <template #default="{ model: proxy, actions: _, save, cancel, isPristine }">
            <TaxonPicker
              label="Taxa"
              v-model="proxy.value.spottings"
              return-object
              :ranks="TaxonRank.ranksUpTo('Family')"
              multiple
              chips
              closable-chips
              clearable
              v-bind="field('spottings')"
            />
            <v-textarea label="Comments" v-model="model.comments" v-bind="field('comments')" />

            <div class="d-flex justify-end">
              <v-btn
                class="mx-1"
                color="primary"
                variant="tonal"
                text="Submit"
                :loading
                @click="submit(proxy.value).then(() => save())"
              />
              <v-btn
                class="mx-1"
                text="Cancel"
                variant="plain"
                color=""
                :disabled="loading"
                @click="(cancel(), toggleEdit(false))"
              />
            </div>
          </template>
        </v-confirm-edit>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { $EventUpdate, Event, EventsService, EventUpdate, TaxonRank } from '@/api'
import { useToggle } from '@vueuse/core'
import TaxonChip from '../taxonomy/TaxonChip.vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import { useSchema } from '../toolkit/forms/schema'

const model = defineModel<Event>({ required: true })

const [editing, toggleEdit] = useToggle(false)
const [loading, toggleLoading] = useToggle(false)

type UpdateData = Pick<EventUpdate, 'spottings' | 'comments'>

const { field, errorHandler } = useSchema($EventUpdate)

async function submit(model: Event) {
  toggleLoading(true)
  const body = toUpdateData(model)
  return EventsService.updateEvent({ path: { id: model.id }, body })
    .then(errorHandler)
    .then(() => toggleEdit(false))
    .finally(() => toggleLoading(false))
}

function toUpdateData(event: Event): UpdateData {
  return {
    spottings: event.spottings?.map(({ name }) => name) ?? null,
    comments: event.comments
  }
}
</script>

<style scoped lang="scss"></style>
