<template>
  <v-card :class>
    <v-row>
      <v-col cols="12">
        <SitePreviewCard :site rounded="0" flat @edit="toggleEdit(true)">
          <template #append>
            <v-btn
              v-if="!!site && !showEdit"
              icon="mdi-pencil"
              size="small"
              variant="tonal"
              @click="toggleEdit(true)"
            />
            <SiteFormDialog
              v-show="!hasID(site) || showEdit"
              v-model:dialog="dialog"
              title="Create site"
              :fullscreen="$vuetify.display.mdAndDown"
              btn-text="Save"
              :max-width="1200"
              @submit="updateSite"
            >
              <template #activator="{ props }">
                <v-btn
                  v-bind="{
                    ...props,
                    ...(!site || 'id' in site
                      ? {
                          text: 'New site',
                          prependIcon: 'mdi-plus'
                        }
                      : {
                          text: 'Edit new site',
                          prependIcon: 'mdi-pencil'
                        })
                  }"
                  variant="tonal"
                  rounded="md"
                >
                </v-btn>
              </template>
            </SiteFormDialog>
          </template>
          <template #default>
            <div v-if="!site || showEdit">
              <v-card-text>
                <SiteAutocomplete @update:model-value="updateSite" />
              </v-card-text>
            </div>
          </template>
          <template #actions v-if="showEdit && site">
            <v-spacer />
            <v-btn text="Cancel" @click="toggleEdit(false)" />
          </template>
        </SitePreviewCard>
      </v-col>
      <v-col cols="12" style="min-height: 400px">
        <SiteProximityMap :model-value="site?.coordinates ?? {}" />
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import SiteFormDialog from '@/components/forms/SiteFormDialog.vue'
import SiteAutocomplete from '@/components/sites/SiteAutocomplete.vue'
import SitePreviewCard from '@/components/sites/SitePreviewCard.vue'
import SiteProximityMap from '@/components/sites/SiteProximityMap.vue'
import { hasID } from '@/functions/db'
import { SiteModel } from '@/models'
import { useToggle } from '@vueuse/core'

const site = defineModel<SiteItem | SiteModel.SiteFormModel>()

defineProps<{
  class?: any
}>()

const dialog = defineModel<boolean>('dialog')
const [showEdit, toggleEdit] = useToggle(false)

function updateSite(s: SiteItem | SiteModel.SiteFormModel | undefined) {
  site.value = s
  dialog.value = false
  toggleEdit(!s)
}
</script>

<style scoped lang="scss"></style>
