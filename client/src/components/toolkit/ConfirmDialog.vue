<template>
  <v-dialog v-model="dialog" max-width="500px">
    <v-card>
      <v-toolbar dark dense flat>
        <v-toolbar-title class="text-body-2 font-weight-bold grey--text">
          {{ content.title }}
        </v-toolbar-title>
      </v-toolbar>
      <v-card-text> {{ content.message }} </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="grey" variant="text" @click="cancel" text="Cancel" />
        <v-btn color="blue-darken-1" variant="text" @click="agree" text="OK" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const dialog = ref(false)
const promise = ref({
  resolve: (_: any) => {},
  reject: (_: any) => {}
})
const content = ref({
  title: '',
  message: ''
})

function open(title: string, message: string) {
  content.value = { title, message }
  dialog.value = true
  return new Promise((resolve, reject) => {
    promise.value = {
      resolve,
      reject
    }
  })
}

function agree() {
  promise.value.resolve(true)
  dialog.value = false
}

function cancel() {
  promise.value.resolve(false)
  dialog.value = false
}

defineExpose({
  open
})
</script>

<style scoped></style>
