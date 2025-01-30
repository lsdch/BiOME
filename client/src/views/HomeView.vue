<template>
  <v-container fluid>
    <v-card v-if="instance" :title="instance.description">
      <template #append>
        <img src="https://img.shields.io/badge/status-in%20development-orange" />
      </template>
    </v-card>
  </v-container>
  <v-sheet min-height="100%" class="d-flex flex-column justify-space-between bg-main fill-height">
    <DashBoard />
  </v-sheet>
  <v-footer border="t-sm">
    <v-container fluid>
      <v-row justify="space-between" align="center">
        <div></div>
        <div>
          <v-btn variant="text"> About </v-btn>
          <v-btn variant="text" :to="{ name: 'api-docs' }"> API Docs </v-btn>
        </div>
        <v-btn
          density="compact"
          stacked
          text="BiOME v0.0"
          variant="text"
          class="font-monospace text-caption"
          color=""
          active-color="primary"
          target="_blank"
          href="https://github.com/lsdch/biome"
        >
          <template #prepend>
            <v-icon size="x-large">mdi-github</v-icon>
          </template>
        </v-btn>
      </v-row>
    </v-container>
  </v-footer>
</template>

<script setup lang="ts">
import { useInstanceSettings } from '@/components/settings'
import DashBoard from '../components/dashboard/DashBoard.vue'

const { instance } = useInstanceSettings()
</script>

<style scoped lang="less"></style>
