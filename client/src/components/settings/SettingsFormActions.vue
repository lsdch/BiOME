<template>
  <v-expand-transition>
    <keep-alive>
      <v-card
        v-if="isOpen"
        class="position-sticky bg-surface d-flex opacity-100 flex-row py-2 align-center"
        scroll-strategy="none"
        :rounded="false"
        style="top: 48px; z-index: 999"
      >
        <slot name="prepend"></slot>
        <v-spacer></v-spacer>
        <v-btn
          color="secondary"
          text="Reset"
          variant="text"
          prepend-icon="mdi-refresh"
          @click="emit('reset')"
        />
        <v-btn
          color="primary"
          text="Save"
          variant="text"
          prepend-icon="mdi-floppy"
          @click="emit('submit')"
        />
      </v-card>
    </keep-alive>
  </v-expand-transition>
</template>

<script setup lang="ts">
const isOpen = defineModel<boolean>()
const emit = defineEmits<{
  reset: []
  submit: []
}>()
</script>

<style scoped></style>
