<template>
  <FormDialog v-model="dialog" :title="`${mode} bio material`" :loading @submit="submit">
    <v-container fluid>
      <v-row v-if="mode === 'Create'">
        <v-col>
          <span class="font-weight-bold mr-3"> Category </span>
          <v-btn-toggle v-model="category" mandatory divided rounded="md" variant="outlined">
            <v-btn
              value="Internal"
              text="Internal"
              :prepend-icon="OccurrenceCategory.icon('Internal')"
              :color="OccurrenceCategory.props.Internal.color"
            />
            <v-btn
              value="External"
              text="External"
              :prepend-icon="OccurrenceCategory.icon('External')"
              :color="OccurrenceCategory.props.External.color"
            />
          </v-btn-toggle>
        </v-col>
      </v-row>
      <v-divider class="my-3" />
      <ExternalBioMatForm v-if="category === 'External'" />
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { BioMaterial, OccurrenceCategory } from '@/api'
import { FormEmits, FormProps, Mode } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { useToggle } from '@vueuse/core'
import { computed, ref } from 'vue'
import ExternalBioMatForm from './ExternalBioMatForm.vue'

const dialog = defineModel<boolean>()

const props = defineProps<FormProps<BioMaterial>>()
const emit = defineEmits<FormEmits<BioMaterial>>()

const mode = computed<Mode>(() => (props.edit ? 'Edit' : 'Create'))
const category = ref<OccurrenceCategory>()

const [loading, toggleLoading] = useToggle(false)

async function submit() {}
</script>

<style scoped lang="scss"></style>
