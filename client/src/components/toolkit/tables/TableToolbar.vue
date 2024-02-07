<template>
  <v-toolbar flat dense prepend-icon="mdi-check" extension-height="auto">
    <template v-if="icon" v-slot:prepend>
      <v-avatar color="secondary" variant="outlined">
        <v-icon dark color="secondary-darken-1">{{ icon }}</v-icon>
      </v-avatar>
    </template>
    <template v-slot:append>
      <slot name="append"></slot>
    </template>

    <v-toolbar-title style="min-width: 150px">{{ title }}</v-toolbar-title>

    <slot v-if="smAndUp && !togglableSearch" name="search" class="flex-grow-1"> </slot>

    <v-spacer />
    <v-btn
      v-if="xs || togglableSearch"
      size="small"
      icon="mdi-magnify"
      color="primary"
      :variant="toggleSearch ? 'flat' : 'text'"
      @click="toggleSearch = !toggleSearch"
    />
    <v-btn
      v-if="form"
      style="min-width: 30px"
      variant="text"
      color="primary"
      :icon="xs"
      size="small"
      @click="formDialog.edit = undefined"
    >
      <v-icon v-if="xs" icon="mdi-plus" size="small" />
      <span v-else>New Item</span>
      <FormDialog
        v-model="formDialog.open"
        activator="parent"
        :entityName="entityName"
        :form="form"
        :edit="formDialog.edit"
        @success="onCreateEdit"
      />
    </v-btn>

    <template v-if="togglableSearch || xs" v-slot:extension>
      <v-expand-transition>
        <div class="w-100 px-3" v-show="toggleSearch" transition="slide-y-transition">
          <slot name="search" class="flex-grow-1">
            <CRUDTableSearchBar v-model="searchTerm" />
          </slot>
        </div>
      </v-expand-transition>
    </template>
    <ConfirmDialog
      v-model="deleteDialog.open"
      v-bind="deleteDialog.props"
      @agree="confirmDeletion"
    />
    <CRUDFeedback v-model="feedback.model" v-bind="feedback.props" />
  </v-toolbar>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import { ApiError, CancelablePromise } from '@/api'
import { HttpStatusCode } from 'axios'
import { Ref, ref } from 'vue'
import { useDisplay } from 'vuetify'
import CRUDFeedback, { FeedbackProps } from '../CRUDFeedback.vue'
import CRUDTableSearchBar from './CRUDTableSearchBar.vue'
import ConfirmDialog, { ConfirmDialogProps } from '../ConfirmDialog.vue'
import FormDialog from '../FormDialog.vue'
import { Mode } from '../form'
import { ToolbarProps } from './table'

const { xs, smAndUp } = useDisplay()

const toggleSearch = ref(false)

const items = defineModel<ItemType[]>({ required: true })

const searchTerm = defineModel<string>('search')

type Props = ToolbarProps<ItemType> & {
  deleteRequest: (item: ItemType) => CancelablePromise<any>
}

const props = defineProps<Props>()

type DeleteDialog<ItemType> = {
  open: boolean
  props: ConfirmDialogProps<ItemType>
}
const deleteDialog: Ref<DeleteDialog<ItemType>> = ref({
  open: false,
  props: {
    title: 'Confirm',
    message: '',
    payload: undefined
  }
})

const formDialog: Ref<{ open: boolean; edit?: ItemType }> = ref({ open: false })

const feedback = ref<{
  model: boolean
  props: FeedbackProps
}>({
  model: false,
  props: {
    text: '',
    color: undefined
  }
})

function showFeedback(text: string, color: string | undefined = undefined) {
  feedback.value.props = { text, color }
  feedback.value.model = true
}

function deleteItem(item: ItemType) {
  deleteDialog.value.props.message = props.itemRepr
    ? `Are you sure you want to delete ${props.itemRepr(item)} ?`
    : 'Are you sure you want to delete this item ?'
  deleteDialog.value.props.payload = item
  deleteDialog.value.open = true
}

async function confirmDeletion(item?: ItemType) {
  if (!item) {
    console.error('Item to delete is undefined. Aborting.')
    return
  }
  console.log(`Deleting item`, item)
  const index = items.value.indexOf(item)
  props
    .deleteRequest(item)
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

function editItem(item: ItemType) {
  formDialog.value.edit = item
  formDialog.value.open = true
}

function onCreateEdit(mode: Mode, item: ItemType) {
  switch (mode) {
    case 'Create':
      console.log('Created item', item)
      items.value.unshift(item)
      showFeedback('Item created', 'success')
      break
    case 'Edit':
      {
        const index = items.value.findIndex(({ id }) => item.id === id)
        if (index < 0) {
          console.error('Failed to find edited item in currently loaded items', item)
          return
        }
        console.log('Edited item', item, ` at index ${index}`)
        items.value.splice(index, 1)
        items.value.unshift(item)
        showFeedback('Item updated', 'success')
      }
      break
    default:
      break
  }
}

defineExpose({ deleteItem, editItem })
</script>

<style scoped></style>
