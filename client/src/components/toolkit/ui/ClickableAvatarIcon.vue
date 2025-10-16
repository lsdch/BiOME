<template>
  <v-tooltip :disabled="!onClick" :content-props="!tooltip ? { class: 'pa-0' } : undefined">
    <template #activator="{ props, isActive }">
      <v-avatar
        variant="outlined"
        :color="isActive ? 'primary' : undefined"
        v-bind="{ ...props, ...$attrs }"
      >
        <v-icon :class="{ 'cursor-default': !onClick }" dark @click="onClick?.()">
          {{ isActive && onClick ? (hoverIcon ?? icon) : icon }}
        </v-icon>
      </v-avatar>
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
defineProps<{
  /**
   * Used to check for event binding
   */
  onClick?: () => void
  icon: string
  hoverIcon?: string
  tooltip?: string
}>()
</script>

<style scoped lang="scss"></style>
