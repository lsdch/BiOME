<template>
  <v-form @submit.prevent="submit">
    <v-row class="mb-3">
      <slot name="prepend-toolbar" :model="model" />
      <v-spacer />
      <v-btn
        color="secondary"
        text="Reset"
        variant="text"
        prepend-icon="mdi-refresh"
        @click="reset"
      />
      <v-btn
        color="primary"
        type="submit"
        text="Save settings"
        variant="text"
        prepend-icon="mdi-floppy"
      />
    </v-row>
    <v-divider />
    <slot :model="model" />
  </v-form>
</template>

<script setup lang="ts" generic="Settings extends {}">
import { CancelablePromise } from '@/api'
import { UnwrapRef, ref, useSlots, watch } from 'vue'

useSlots()

const props = defineProps<{
  get(): CancelablePromise<Settings>
  update(data: { requestBody: UnwrapRef<Awaited<Settings>> }): CancelablePromise<Settings>
}>()

const initial = await props.get()

const model = ref(initial)
const disabled = ref(true)

watch(model, (m) => {
  disabled.value = m == initial
})

function submit() {
  const resp = props.update({ requestBody: model.value })
  resp.then((_settings) => {
    console.log('Settings saved')
  })
}

function reset() {
  model.value = initial as UnwrapRef<Awaited<Settings>>
}
</script>

<style scoped></style>
