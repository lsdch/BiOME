<!-- eslint-disable vue/valid-v-slot -->
<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-stepper editable :items="stepper_items">
          <template v-slot:item.1>
            <v-text-field
              class="mt-3"
              name="institution"
              label="Institution name"
              id="institution-input"
              hint="The name of the structure for which this instance is deployed, e.g. the name of your lab."
              persistent-hint
            />
            <v-switch
              label="Allow public access"
              hint="Anonymous users will have read-only access to all non-personal data on the platform."
              persistent-hint
            />
          </template>
          <template v-slot:item.2>
            <AnchorTaxonPicker />
          </template>
        </v-stepper>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import AnchorTaxonPicker from '@/components/taxonomy/imports/AnchorTaxonPicker.vue'
const stepper_items = ['General configuration', 'Admin account', 'Backbone taxonomy']
</script>

<style scoped></style>
