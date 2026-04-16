<template>
  <div>
    <v-combobox
      class="required"
      label="Collection"
      v-model.trim="model.name"
      :items="collections"
      density="compact"
      v-bind="schema('name')"
      clearable
    >
      <template #selection="{ item, index }">
        <div class="d-flex ga-3">
          <v-chip
            v-if="item && !collections?.includes(item)"
            size="small"
            color="success"
            text="New"
          />
          <span>{{ item }}</span>
        </div>
      </template>
    </v-combobox>
    <v-combobox
      label="Vouchers"
      multiple
      v-model.trim="model.vouchers"
      density="compact"
      v-bind="schema('vouchers')"
      clearable
      chips
      closable-chips
    ></v-combobox>
  </div>
</template>

<script setup lang="ts">
import { $CollectionField, CollectionField } from '@/api'
import { listCollectionsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useSchema } from '@/composables/schema'
import { useQuery } from '@tanstack/vue-query'

const { data: collections } = useQuery(listCollectionsOptions())

const model = defineModel<CollectionField>({ required: true })

const {
  bind: { schema }
} = useSchema($CollectionField)
</script>

<style scoped lang="scss"></style>
