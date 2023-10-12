<template>
  <div id="swagger-ui" />
</template>

<script setup lang="ts">
import SwaggerUI from 'swagger-ui'
import 'swagger-ui/dist/swagger-ui.css'
import { onMounted } from 'vue'

onMounted(() => {
  SwaggerUI({
    url: '/api/v1/swagger/doc.json',
    dom_id: '#swagger-ui'
  })
})
</script>

<style scoped></style>
