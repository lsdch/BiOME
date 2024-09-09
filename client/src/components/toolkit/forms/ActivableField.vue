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
import { UserRole } from '@/api'
import { useUserStore } from '@/stores/user'
import { computed, ref } from 'vue'
import { VConfirmEdit } from 'vuetify/components'

const confirmEdit = ref<VConfirmEdit>()

const model = defineModel<Value>()

const props = defineProps<{ activable?: true | UserRole }>()

const active = ref(false)

const { isGranted } = useUserStore()

function save() {
  confirmEdit.value!.save()
}

function cancel() {
  confirmEdit.value!.cancel()
}

const isActivable = computed(() => {
  return props.activable === true || (!!props.activable && isGranted(props.activable))
})

const slotProps = computed(() => {
  return {
    onfocus() {
      active.value = isActivable.value
    },
    readonly: !active.value
  }
})
</script>

<style scoped></style>
