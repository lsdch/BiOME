<template>
  <v-confirm-edit ref="confirmEdit" v-model="model" @save="active = false" @cancel="active = false">
    <template #default="{ model: proxy, actions, isPristine }">
      <div @focusout="isPristine ? (active = false) : null">
        <slot :proxy :save :cancel :actions :isPristine :active :props="slotProps" />
      </div>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts" generic="Value">
import { computed, ref } from 'vue'
import { VConfirmEdit } from 'vuetify/components'

const confirmEdit = ref<VConfirmEdit>()

const model = defineModel<Value>()

const active = ref(false)

function save() {
  confirmEdit.value!.save()
}

function cancel() {
  confirmEdit.value!.cancel()
}

const slotProps = computed(() => {
  return {
    onfocus() {
      active.value = true
    },
    readonly: !active.value
  }
})
</script>

<style scoped></style>
