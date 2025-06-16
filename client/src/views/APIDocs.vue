<template>
  <elements-api
    apiDescriptionUrl="/api/v1/openapi.json"
    router="hash"
    layout="responsive"
    tryItCredentialsPolicy="same-origin"
  />
</template>

<script setup lang="ts">
import '@stoplight/elements/web-components.min.js'
import '@stoplight/elements/styles.min.css'
</script>

<style lang="scss">
html[data-theme='dark'] {
  .sl-code-viewer {
    background: #fafafa;
    color: black;
    filter: invert(1);
  }
}
</style>
