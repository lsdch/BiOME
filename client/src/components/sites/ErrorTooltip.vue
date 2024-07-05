<template>
  <v-tooltip v-if="error">
    <template #activator="{ props: p }">
      <div v-bind="p" :class="[props.class, { [props.errorClass]: error }]">
        <slot> </slot>
      </div>
    </template>
    {{ error }}
  </v-tooltip>
  <div v-else :class="props.class">
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    error?: string
    class?: string
    errorClass?: string
  }>(),
  {
    errorClass: 'text-error'
  }
)
</script>

<style scoped></style>
