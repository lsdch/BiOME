<template>
  <v-stepper v-model="step" class="fill-height d-flex flex-column" :mobile="mobile">
    <v-stepper-header>
      <div class="d-flex align-center">
        <v-btn
          class="mx-2"
          color="secondary"
          icon="mdi-arrow-left"
          variant="text"
          :to="{ name: 'sites' }"
        />
        <v-stepper-item :value="1" title="Dataset" editable />
      </div>
      <v-stepper-item :value="2" title="New sites" editable />
      <v-stepper-item :value="3" title="Review" />
    </v-stepper-header>

    <v-stepper-window class="fill-height">
      <v-stepper-window-item :value="1">
        <v-container>
          <SiteDatasetImport />
        </v-container>
      </v-stepper-window-item>
      <v-stepper-window-item :value="2">
        <SiteTabularImport />
      </v-stepper-window-item>
      <v-stepper-window-item :value="3" class="fill-height">
        <div class="fill-height w-100">
          <SitesMap />
        </div>
      </v-stepper-window-item>
    </v-stepper-window>

    <v-stepper-actions @click:next="step += 1" @click:prev="step -= 1" />
  </v-stepper>
</template>

<script setup lang="ts">
import SitesMap from '@/components/maps/SitesMap.vue'
import SiteDatasetImport from '@/components/sites/SiteDatasetImport.vue'
import SiteTabularImport from '@/components/sites/SiteTabularImport.vue'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'

const { mobile } = useDisplay()

const step = ref(1)
</script>

<style lang="scss">
.v-stepper-window {
  margin: 0;
}
</style>
