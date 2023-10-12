<script setup lang="ts">
import DashBoard from '../components/dashboard/DashBoard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <v-container class="dashboard">
    <DashBoard />

    <v-footer>
      <v-row justify="center">
        <v-btn variant="text"> About </v-btn>
        <v-btn variant="text" :to="router.resolve({ name: 'api-docs' })"> API Docs </v-btn>
      </v-row>
    </v-footer>
  </v-container>
</template>

<style scoped lang="less">
.dashboard {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  justify-content: space-between;
}
</style>
