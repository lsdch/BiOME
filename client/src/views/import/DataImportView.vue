<template>
  <v-container class="fill-height">
    <div class="import-cards-container d-flex flex-column align-center w-100">
      <v-card title="Sites" prepend-icon="mdi-map-marker">
        <v-divider />
        <v-list>
          <SiteFormDialog>
            <template #activator="{ props, isActive }">
              <v-list-item
                title="Create record"
                prepend-icon="mdi-plus-circle-outline"
                v-bind="props"
                :active="isActive"
                color="primary"
              />
            </template>
          </SiteFormDialog>
          <v-list-item prepend-icon="mdi-file-upload-outline">Bulk import</v-list-item>
        </v-list>
      </v-card>
      <v-card title="Occurrences" prepend-icon="mdi-crosshairs-gps">
        <v-divider />
        <!-- <v-card-text class="d-flex justify-space-around">
            <v-btn
              stacked
              text="Create record"
              prepend-icon="mdi-plus-circle-outline"
              variant="outlined"
              rounded="md"
              class="w-50"
            />
            <v-btn
              stacked
              text="Bulk import"
              prepend-icon="mdi-file-upload-outline"
              variant="outlined"
              rounded="md"
              class="w-50"
            />
          </v-card-text> -->
        <v-list>
          <BioMaterialFormDialog>
            <template #activator="{ props, isActive }">
              <v-list-item
                title="Create record"
                prepend-icon="mdi-plus-circle-outline"
                v-bind="props"
                :active="isActive"
                color="primary"
              />
            </template>
          </BioMaterialFormDialog>
          <!-- <v-list-item prepend-icon="mdi-plus-circle-outline">Create record</v-list-item> -->
          <v-list-item prepend-icon="mdi-file-upload-outline">Bulk import</v-list-item>
        </v-list>
      </v-card>

      <v-card title="Sequences" prepend-icon="mdi-dna">
        <v-divider />
        <v-list>
          <v-list-item prepend-icon="mdi-plus-circle-outline">Create record</v-list-item>
          <v-list-item prepend-icon="mdi-file-upload-outline">Bulk import</v-list-item>
        </v-list>
      </v-card>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import BioMaterialFormDialog from '@/components/occurrence/BioMaterialFormDialog.vue'
import SiteFormDialog from '@/components/sites/SiteFormDialog.vue'
</script>

<style scoped lang="scss">
.import-cards-container {
  gap: 2rem;
  .v-card {
    width: 100%;
    max-width: 500px;
  }
}
</style>
