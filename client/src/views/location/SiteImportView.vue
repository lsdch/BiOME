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
        <v-stepper-item
          :value="1"
          title="Dataset"
          editable
          :rules="[() => validity.step1 ?? true]"
        />
      </div>
      <v-stepper-item :value="2" title="Sites coordinates" />
      <v-stepper-item :value="3" title="New sites" editable />
      <v-stepper-item :value="4" title="Review" />
    </v-stepper-header>

    <v-stepper-window class="fill-height">
      <v-stepper-window-item :value="1">
        <v-container>
          <v-form v-model="validity.step1">
            <v-row>
              <v-col>
                <v-text-field v-model="model.label" v-bind="field('label')" label="Dataset name" />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <PersonPicker
                  v-bind="field('maintainers')"
                  @update:modelValue="
                    (v) => {
                      if (Array.isArray(v)) {
                        model.maintainers = v.map(({ alias }) => alias)
                      } else {
                        model.maintainers = [v.alias]
                      }
                    }
                  "
                  label="Maintainers"
                  multiple
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-textarea
                  v-model="model.description"
                  v-bind="field('description')"
                  label="Description"
                  variant="outlined"
                />
              </v-col>
            </v-row>
          </v-form>
        </v-container>
      </v-stepper-window-item>
      <v-stepper-window-item :value="2" class="fill-height">
        <SiteDatasetPrimer />
      </v-stepper-window-item>
      <v-stepper-window-item :value="3">
        <SiteTabularImport />
      </v-stepper-window-item>
      <v-stepper-window-item :value="4" class="fill-height">
        <div class="fill-height w-100">
          <SitesMap />
        </div>
      </v-stepper-window-item>
    </v-stepper-window>

    <v-stepper-actions @click:next="step += 1" @click:prev="step -= 1" />
  </v-stepper>
</template>

<script setup lang="ts">
import { $SiteDatasetInput, SiteDatasetInput } from '@/api'
import SitesMap from '@/components/maps/SitesMap.vue'
import PersonPicker from '@/components/people/PersonPicker.vue'
import SiteDatasetPrimer from '@/components/sites/SiteDatasetPrimer.vue'
import SiteTabularImport from '@/components/sites/SiteTabularImport.vue'
import { FormProps, useForm } from '@/components/toolkit/forms/form'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'

const { mobile } = useDisplay()

const step = ref(1)
const validity = ref({
  step1: undefined
})

const props = defineProps<FormProps<SiteDatasetInput>>()

const initial: SiteDatasetInput = {
  label: '',
  maintainers: []
}
const { field, model } = useForm(props, $SiteDatasetInput, { initial })
</script>

<style lang="scss">
.v-stepper-window {
  margin: 0;
}
</style>
