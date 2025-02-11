<template>
  <v-card v-bind="$attrs">
    <slot></slot>
    <template #append>
      <slot name="append"></slot>
    </template>
  </v-card>
  <CardDialog v-model="dialog" v-bind="$attrs">
    <template #append>
      <slot name="append"></slot>
    </template>
    <v-card-text class="fill-height">
      <slot></slot>
    </v-card-text>
  </CardDialog>
</template>

<script setup lang="ts">
import CardDialog from '../forms/CardDialog.vue'

const dialog = defineModel<boolean>({ default: false })
</script>

<style scoped lang="scss"></style>
