<script setup lang="ts">
import DashBoard from '../components/dashboard/DashBoard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="d-flex flex-column justify-space-between fill-height">
    <DashBoard />
    <v-footer class="flex-grow-0">
      <v-row justify="center">
        <v-btn variant="text"> About </v-btn>
        <v-btn variant="text" :to="router.resolve({ name: 'api-docs' })"> API Docs </v-btn>
      </v-row>
    </v-footer>
  </div>
</template>

<style scoped lang="less"></style>
