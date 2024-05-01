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
    <slot :model="model" :schema="schema" />
  </v-form>
</template>

<script setup lang="ts" generic="Settings extends {}, P extends SchemaProperties">
import { Ref, computed, ref, useSlots, watch } from 'vue'
import { SchemaProperties, SchemaWithProperties, useSchema } from '../toolkit/form'

useSlots()

const props = defineProps<{
  get(): PromiseLike<Settings>
  update(data: { requestBody: Settings | Awaited<Settings> }): PromiseLike<Settings>
  schema?: SchemaWithProperties<P>
}>()

const initial = await props.get()

const model = ref<Awaited<Settings>>(initial) as Ref<Settings>
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
  model.value = initial
}

const schema = computed(() => {
  if (props.schema !== undefined) {
    const { schema: inputSchema } = useSchema(props.schema)
    return inputSchema
  }
  return () => ({})
})
</script>

<style scoped></style>
