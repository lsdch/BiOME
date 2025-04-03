<template>
  <v-select :items="$Quantity.enum" v-model="model" v-bind="$attrs">
    <template #item="{ props, item }">
      <v-list-item v-bind="props">
        <template #append>
          <Quantity.Chip :quantity="item.raw" size="small" class="font-monospace" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { $Quantity, Quantity } from '@/api'

const model = defineModel<Quantity>()
</script>

<style scoped lang="scss"></style>
