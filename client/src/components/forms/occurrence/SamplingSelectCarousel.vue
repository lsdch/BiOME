<template>
  <template v-if="samplings.length">
    <v-carousel
      v-model="selection"
      @update:model-value="(v) => (model = samplings[v])"
      color="primary"
      class="sampling-carousel"
      height="auto"
      hide-delimiter-background
      content-class="justify-center ga-2"
      :show-arrows="samplings.length > 1"
    >
      <template #prev="{ props }">
        <v-btn v-bind="props" size="small" :rounded="100" variant="tonal"></v-btn>
      </template>
      <template #next="{ props }">
        <v-btn v-bind="props" size="small" :rounded="100" variant="tonal"></v-btn>
      </template>
      <v-carousel-item v-for="(sampling, i) in samplings">
        <v-list class="pb-12">
          <SamplingListItems :sampling />
        </v-list>
      </v-carousel-item>
    </v-carousel>
  </template>
</template>

<script setup lang="ts" generic="Item extends Sampling | SamplingInner">
import { Sampling, SamplingInner } from '@/api'
import SamplingListItems from '@/components/events/SamplingListItems.vue'
import { ref } from 'vue'

defineProps<{ samplings: Item[] }>()
const selection = ref<number>(0)

const model = defineModel<Item>()
</script>

<style lang="scss">
.sampling-carousel {
  // .v-window__container {
  //   background-color: transparent;
  //   z-index: 1000;
  // }
  .v-window__controls {
    align-items: end;
    padding-bottom: 5px;
    button {
      z-index: 1000;
    }
    // position: relative;
  }
}
</style>
