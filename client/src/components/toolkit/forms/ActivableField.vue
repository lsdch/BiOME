<template>
  <v-confirm-edit ref="confirmEdit" v-model="model" @save="active = false" @cancel="active = false">
    <template #default="{ model: proxy, actions, isPristine }">
      <div class="d-flex align-center" @focusout="isPristine ? (active = false) : null">
        <slot :proxy :save :cancel :actions :isPristine :active :props="slotProps" />
        <v-btn
          v-show="active && !isPristine"
          class="flex-grow-0"
          color="success"
          icon="mdi-check"
          density="compact"
          variant="plain"
          @click="save()"
        />
        <v-btn
          v-show="active && !isPristine"
          class="flex-grow-0"
          color="error"
          icon="mdi-close"
          density="compact"
          variant="plain"
          @click="cancel()"
        />
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
