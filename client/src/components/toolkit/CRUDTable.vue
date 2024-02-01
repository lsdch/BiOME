<template>
  <div>
    <v-data-table :headers="headers" :items="items" :loading="loading" v-bind="$attrs">
      <template v-slot:top>
        <v-toolbar flat prepend-icon="mdi-check">
          <template v-if="icon" v-slot:prepend>
            <v-icon :icon="icon" />
          </template>
          <v-toolbar-title :text="title" />
          <v-spacer />

          <v-btn variant="text" color="primary" @click="editItem = undefined">
            New Item
            <FormDialog
              ref="formDialog"
              activator="parent"
              :form="form"
              :edit="editItem"
              @success="onSuccess"
            />
          </v-btn>
          <ConfirmDialog ref="deleteDialog" />
        </v-toolbar>
      </template>
      <template v-if="props.showActions" v-slot:[`item.actions`]="{ item }">
        <v-icon size="small" icon="mdi-pencil" @click="modifyItem(item)" />
        <v-icon size="small" icon="mdi-delete" @click="deleteItem(item)" />
      </template>
      <template v-for="(name, index) of slotNames" v-slot:[name]="slotData" :key="index">
        <slot :name="name" v-bind="slotData || {}" />
      </template>
    </v-data-table>
    <v-snackbar v-model="feedback.show" :timeout="2000" :color="feedback.color">
      {{ feedback.text }}
      <template v-slot:actions>
        <v-btn variant="text" @click="feedback.show = false"> Close </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup lang="ts" generic="ItemInputType extends {}, ItemType extends { id: string }">
import { useSlots, type Component } from 'vue'
import { ApiError, CancelablePromise } from '@/api'
import { Ref, ref, onMounted } from 'vue'
import { HttpStatusCode } from 'axios'
import { type VDataTable, VSnackbar } from 'vuetify/components'
import ConfirmDialog from './ConfirmDialog.vue'
import FormDialog from './FormDialog.vue'
import type { ComponentExposed } from 'vue-component-type-helpers'
import { Prop } from 'vue'
import { Mode } from './form'

const slots = useSlots()
// Assert type here to prevent errors in template
const slotNames = Object.keys(slots) as 'default'[]

const loading = ref(true)
const formDialog = ref<ComponentExposed<typeof FormDialog> | null>(null)
const deleteDialog = ref<ComponentExposed<typeof ConfirmDialog> | null>(null)
const editItem: Ref<ItemType | undefined> = ref(undefined)

type Props = {
  title: string
  headers: ReadonlyHeaders
  list: () => CancelablePromise<ItemType[]>
  delete: (item: ItemType) => CancelablePromise<any>
  create: (item: ItemInputType) => CancelablePromise<any>
  update: (item: ItemType) => CancelablePromise<any>
  form: Component<Prop<{ onSuccess: string }>>
  showActions?: boolean
  itemRepr?: (item: ItemType) => string
  icon?: string
}

const props = defineProps<Props>()
defineSlots<VDataTable['$slots']>()

const items: Ref<ItemType[]> = ref([])

onMounted(async () => {
  items.value = await props.list()
  loading.value = false
})

const feedback = ref<{
  show: boolean
  text: string
  color: string | undefined
}>({
  show: false,
  text: '',
  color: undefined
})
function showFeedback(text: string, color: string | undefined = undefined) {
  feedback.value = {
    show: true,
    text,
    color
  }
}

function onSuccess(mode: Mode, item: ItemType) {
  switch (mode) {
    case 'Create':
      items.value.unshift(item)
      showFeedback('Item created', 'success')
      break
    case 'Edit':
      console.log('Edited item', item)
      items.value = items.value.filter((oldItem) => oldItem.id !== item.id)
      items.value.unshift(item)
      showFeedback('Item updated', 'success')
      break
    default:
      break
  }
}

function modifyItem(item: ItemType) {
  editItem.value = item
  formDialog.value?.open()
}

async function deleteItem(item: ItemType) {
  let msg = props.itemRepr
    ? `Are you sure you want to delete ${props.itemRepr(item)} ?`
    : 'Are you sure you want to delete this item ?'
  await deleteDialog.value?.open('Confirm', msg).then((confirm) => {
    if (confirm === true) {
      console.log(`Deleting item`, item)
      let index = items.value.indexOf(item)
      props
        .delete(item)
        .then(() => {
          items.value.splice(index, 1)
          showFeedback('Item successfully deleted.', 'success')
        })
        .catch((err: ApiError) => {
          switch (err.status) {
            case HttpStatusCode.NotFound:
              showFeedback('Deletion failed: record not found.', 'error')
              break
            case HttpStatusCode.BadRequest:
              showFeedback(`Deletion was not allowed: ${err.message}`, 'error')
              break
            case HttpStatusCode.Forbidden:
              showFeedback('You are not granted rights to modify this item.', 'error')
              break
            case HttpStatusCode.InternalServerError:
              showFeedback('An unexpected error occurred.', 'error')
          }
        })
    }
  })
}
</script>

<style scoped></style>
