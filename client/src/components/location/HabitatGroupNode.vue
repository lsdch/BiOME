<template>
  <HabitatNodeContainer
    :sideLabel="data.elements.length > 1 ? data.label : undefined"
    layout="vertical"
    :selected="selected"
  >
    <HabitatElement v-for="child in data.elements" :key="child.id" :habitat="child" :group="data" />
  </HabitatNodeContainer>
</template>

<script setup lang="ts">
import { HabitatGroup } from '@/api'
import { type NodeProps } from '@vue-flow/core'
import HabitatElement from './HabitatElement.vue'
import HabitatNodeContainer from './HabitatNodeContainer.vue'

const props = defineProps<NodeProps<HabitatGroup>>()

// const isRootGroup = computed(() => !props.data.depends)
</script>

<style lang="scss"></style>
