<template>
  <div>
    <v-alert v-if="error" color="error">
      Failed to load samplings for site {{ props.siteCode }}: {{ error }}
    </v-alert>
    <slot name="no-data" v-else-if="!isFetching && !items?.length" />
    <v-autocomplete
      v-else
      v-model="model"
      :loading="isFetching"
      :items
      item-value="id"
      return-object
    >
      <template #selection="{ item: sampling }">
        <v-list-item :title="DateWithPrecision.format(sampling.performed_on)" class="px-0">
          <template #append>
            <v-chip :text="sampling.number.toString()" size="small" prepend-icon="mdi-pound" />
          </template>
        </v-list-item>
      </template>

      <template #item="{ item: sampling, props }">
        <v-list-item v-bind="props" :title="DateWithPrecision.format(sampling.performed_on)">
          <template #append>
            <v-chip :text="sampling.number.toString()" size="small" prepend-icon="mdi-pound" />
          </template>
        </v-list-item>
      </template>
    </v-autocomplete>
  </div>
</template>

<script setup lang="ts">
import { DateWithPrecision, Sampling } from '@/api'
import { listSiteSamplingsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const model = defineModel<Sampling>()

const props = defineProps<{ siteCode: string }>()

const {
  data: items,
  isFetching,
  error
} = useQuery(
  computed(() => ({
    enabled: !!props.siteCode,
    ...listSiteSamplingsOptions({
      path: { code: props.siteCode }
    })
  }))
)
</script>

<style scoped lang="scss"></style>
