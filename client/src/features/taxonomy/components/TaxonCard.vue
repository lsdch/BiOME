<template>
  <v-bottom-sheet v-model="open" :inset="mdAndUp" content-class="rounded-0">
    <v-card :rounded="false" :title="taxon.name" :subtitle="taxon.authorship" :loading="isFetching">
      <template #prepend>
        <LinkIconGBIF v-if="taxon.gbif_id" :GBIF_ID="taxon.gbif_id" variant="text" />
        <FTaxonStatusIndicator v-else :status="taxon.status" />
      </template>

      <template #append>
        <v-btn
          v-if="taxon.status != 'accepted'"
          variant="plain"
          color="primary"
          text="EDIT"
          prepend-icon="mdi-pencil-outline"
        />
        <v-btn variant="text" icon="mdi-close" @click="open = false" />
      </template>

      <template #text>
        <v-chip :text="taxon.rank" variant="outlined" class="mr-3" />
        <v-chip
          :text="taxon.status"
          variant="outlined"
          :color="taxonStatusIndicatorProps(taxon.status).color"
        />
        <!-- <v-col cols="12" sm="6">
            <ActivableField v-model="taxon.code" activable="Maintainer">
              <template #default="{ proxy, active, props, save, cancel, isPristine }">
                <v-text-field
                  label="Code"
                  v-model="proxy.value"
                  v-bind="{ ...props, ...schema('code') }"
                  :variant="active ? 'underlined' : 'plain'"
                >
                  <template #append="{ isValid }">
                    <v-btn
                      v-if="active && !isPristine"
                      :disabled="!isValid.value"
                      class="flex-grow-0"
                      color="success"
                      icon="mdi-check"
                      density="compact"
                      variant="plain"
                      @click="save"
                    />
                    <v-btn
                      v-if="active && !isPristine"
                      class="flex-grow-0"
                      color="error"
                      icon="mdi-close"
                      density="compact"
                      variant="plain"
                      @click="cancel"
                    />
                  </template>
</v-text-field>
</template>
</ActivableField>
</v-col> -->

        <div>
          <ActivableField
            v-model="taxon.comments"
            v-if="taxon.comments || isGranted('maintainer')"
            activable="maintainer"
          >
            <template #default="{ proxy, active, props, actions, isPristine }">
              <v-textarea
                v-model="proxy.value"
                :label="proxy.value || active ? 'Comment' : 'Add comment...'"
                :rows="1"
                auto-grow
                v-bind="{ ...props, ...schema('comments') }"
                :variant="active ? 'underlined' : 'plain'"
              >
                <template #details>
                  <div class="align-self-start" v-if="active && !isPristine">
                    <component :is="actions"></component>
                  </div>
                </template>
              </v-textarea>
            </template>
          </ActivableField>
        </div>

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
                :title="v?.rank"
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
          <v-chip
            color="primary"
            :text="`${relatives?.descendants.length}`"
            :rounded="100"
            size="small"
          />
        </v-list-subheader>
        <div class="descendants">
          <v-alert v-if="error" color="error"> Failed to retrieve descendants list </v-alert>
          <v-skeleton-loader type="chip@5">
            <v-chip
              v-for="c in relatives?.descendants"
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
        <!-- <div>
          <ItemDateChip v-if="taxon.meta?.created" icon="created" :date="taxon.meta.created" />
          <ItemDateChip v-if="taxon.meta?.modified" icon="updated" :date="taxon.meta.modified" />
        </div> -->
        <v-spacer />
        <div v-if="isGranted('admin')">
          <DeleteBtn
            :title="`Confirm deletion : ${taxon.name}`"
            @confirm="
              deleteTaxon({
                path: { id: taxon.id }
              })
            "
          >
            <template #message>
              <v-card-text>
                <b> Are you sure you want to delete this taxon and all of its descendants ? </b>
                <ul>
                  <li>All descendants of this taxon will be deleted.</li>
                  <li>
                    All occurrences associated with this taxon and its descendants will be deleted.
                  </li>
                  <li>
                    Samplings associated with deleted occurrences will be deleted, if no other
                    occurrences are associated with them.
                  </li>
                </ul>
                <p class="text-error">This action cannot be undone.</p>
                <v-text-field
                  label="Please type the exact taxon name to confirm deletion"
                  class="required"
                  :rules="[(v) => v === taxon.name || 'Name does not match']"
                ></v-text-field>
              </v-card-text>
            </template>
          </DeleteBtn>
          <v-btn
            v-if="TaxonRank.extensibleRanks.includes(taxon.rank)"
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
import { $CreateTaxonInput, Taxon, TaxonRank } from '@/api'
import { deleteTaxonMutation, getTaxonOptions } from '@/api/gen/@tanstack/vue-query.gen'
import ActivableField from '@/components/toolkit/forms/ActivableField.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import { useAppConfirmDialog } from '@/composables/confirm_dialog'
import { useSchemaBinding } from '@/composables/schema'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import { useDisplay } from 'vuetify'
import { FTaxonStatusIndicator, taxonStatusIndicatorProps } from './functionals'
import LinkIconGBIF from './LinkIconGBIF'

const { mdAndUp } = useDisplay()
const { isGranted } = useUserStore()

const taxon = defineModel<Taxon>({ required: true })
const open = defineModel<boolean>('open')

const { schema } = useSchemaBinding($CreateTaxonInput)

const emit = defineEmits<{
  'add-child': [parent: Taxon]
  deleted: []
  navigate: [target: Taxon]
}>()

const {
  data: relatives,
  error,
  isFetching
} = useQuery(computed(() => getTaxonOptions({ path: { id: taxon.value.id } })))

const { askConfirm } = useAppConfirmDialog()
const { feedback } = useFeedback()

const { mutate: deleteTaxon, isPending: isDeleting } = useMutation({
  ...deleteTaxonMutation(),
  onSuccess: () => {
    emit('deleted')
    feedback({
      type: 'success',
      message: `Taxon successfully deleted along with all of its descendants`
    })
    open.value = false
  },
  onError: (error) => {
    feedback({ type: 'error', message: 'Failed to delete taxon' })
    console.error(error)
  }
})

// async function deleteTaxon(taxon: Taxon) {
//   askConfirm({
//     title: `Delete taxon ${taxon.name}?`,
//     message: 'All descendants will also be deleted'
//   }).then(async ({ isCanceled }) => {
//     if (isCanceled) return
//     delTaxon({ path: { name: taxon.name } })
//   })
// }
</script>

<style scoped>
.descendants {
  max-height: 50dvh;
  overflow-y: scroll;
}
</style>
