<template>
  <Handle :position="Position.Top" type="target" style="visibility: hidden" />
  <Handle :position="Position.Bottom" type="source" />
  <HabitatNodeContainer :selected="selected">
    <HabitatElement :habitat="data"> </HabitatElement>
  </HabitatNodeContainer>
</template>

<script setup lang="ts">
import { Habitat } from '@/api'
import HabitatNodeContainer from './HabitatNodeContainer.vue'
import HabitatElement from './HabitatElement.vue'
import { Handle, NodeProps, Position } from '@vue-flow/core'

defineProps<NodeProps<Habitat>>()
</script>

<style scoped></style>
