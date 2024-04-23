<template>
  <elements-api
    apiDescriptionUrl="/api/v1/openapi.json"
    router="hash"
    layout="responsive"
    basePath="/docs/api"
    tryItCredentialsPolicy="same-origin"
  />
</template>

<script setup lang="ts">
import '@stoplight/elements/web-components.min.js'
import '@stoplight/elements/styles.min.css'
</script>

<style scoped></style>
