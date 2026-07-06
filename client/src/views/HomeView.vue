<template>
  <v-container fluid>
    <v-card v-if="instance" :title="instance.subtitle" :subtitle="instance.description">
      <template #append>
        <img src="https://img.shields.io/badge/status-in%20development-orange" />
      </template>
    </v-card>
  </v-container>
  <v-sheet min-height="100%" class="d-flex flex-column justify-space-between bg-main">
    <DashBoard />
  </v-sheet>
  <v-footer border="t-sm">
    <v-container fluid>
      <v-row class="justify-space-between align-center">
        <div></div>
        <div>
          <v-btn variant="text"> About </v-btn>
          <v-btn variant="text" :to="{ name: 'api-docs' }"> API Docs </v-btn>
        </div>
        <div class="d-flex align-center gap-2">
          <v-btn
            variant="text"
            class="font-monospace text-caption"
            color=""
            rounded="md"
            active-color="primary"
            target="_blank"
            href="https://github.com/lsdch/biome"
          >
            <div class="d-flex align-end">
              <img
                src="@/../../biome_logo.svg"
                alt="BiOME Logo"
                height="24"
                class="py-1"
                style="filter: brightness(0) invert(1)"
              />
              <span>v0.0</span>
            </div>
            <template #prepend>
              <v-icon size="x-large">mdi-github</v-icon>
            </template>
          </v-btn>
        </div>
      </v-row>
    </v-container>
  </v-footer>
</template>

<script setup lang="ts">
import { useInstanceSettings } from '@/features/settings/components'
import DashBoard from '@/features/dashboard/views/DashBoard.vue'

const { instance } = useInstanceSettings()
</script>

<style scoped lang="less"></style>
