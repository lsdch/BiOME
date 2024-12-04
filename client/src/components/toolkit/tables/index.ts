import { ErrorModel } from "@/api"
import { useAppConfirmDialog } from "@/composables/confirm_dialog"
import { useUserStore } from "@/stores/user"
import { RequestResult } from "@hey-api/client-fetch"
import { HttpStatusCode } from "axios"
import { MaybeRef, ModelRef, computed, onMounted, ref, triggerRef } from "vue"
import { FeedbackProps } from "../CRUDFeedback.vue"
import { Mode } from "../forms/form"



export type ToolbarProps = {
  /**
   * Table toolbar title
   */
  title?: string
  /**
   * Table icon
   */
  icon?: string
  /**
   * Whether table search bar is displayed by default, or must be toggled
   */
  togglableSearch?: boolean
  /**
   * Disable sorting by item last update
   */
  noSort?: boolean
  /**
   * Disable common filters such as owned items
   */
  noFilters?: boolean
  /**
   * Used to check whether a reload event listener is bound to the toolbar
   */
  onReload?: Function
}

export type TableProps<ItemType> = {
  /**
   * Entity name to display as title
   */
  entityName: string
  /**
   * Datatable headers definition
   */
  headers: CRUDTableHeader[]
  /**
   * Table toolbar configuration
   */
  toolbar?: ToolbarProps | false
  /**
   * Display column with update/delete controls
   */
  appendActions?: boolean | "delete" | "edit"
  /**
   * Short string representation of a table item to display in UI
   */
  itemRepr?: (item: ItemType) => string
  /**
   * API call to populate table items
   */
  fetchItems?: () => RequestResult<ItemType[], ErrorModel, false>
  /**
   * API call to delete an item
   */
  delete?: (item: ItemType) => RequestResult<ItemType, ErrorModel, false>
  /**
   * Reload all items after deleting one
   */
  reloadOnDelete?: boolean
}

export type TableEmits<ItemType> = (
  ((evt: "itemCreated", item: ItemType, index: number) => void) &
  ((evt: "itemEdited", item: ItemType, index: number) => void)
)

type FormSlotScope<ItemType extends { id: string }> = {
  dialog: boolean,
  mode: Mode,
  editItem?: MaybeRef<ItemType>,
  onSuccess: (_item: ItemType) => void,
  onClose: () => void,
}


export function useTable<ItemType extends { id: string }>(
  items: ModelRef<ItemType[]>,
  props: TableProps<ItemType>,
  emit: TableEmits<ItemType>
) {

  const { user: currentUser } = useUserStore()
  const { askConfirm } = useAppConfirmDialog()

  const form = ref<FormSlotScope<ItemType>>({
    dialog: false,
    mode: 'Create',
    editItem: undefined,
    onSuccess: (_item: ItemType) => { },
    onClose: () => { }
  })

  const processedHeaders = computed((): CRUDTableHeader[] => {
    const headersWithActions = props.appendActions && currentUser !== undefined && currentUser.role !== "Visitor"
      ? props.headers.concat([{ title: 'Actions', key: 'actions', sortable: false, align: 'center', width: 0, cellProps: { class: 'text-no-wrap' } }])
      : props.headers
    return headersWithActions.filter(({ hide }) => {
      return !hide?.value
    })
  })

  const loading = ref(props.fetchItems !== undefined)
  const loadingFailed = ref(false)
  async function loadItems() {
    if (props.fetchItems) {
      loading.value = true
      const { data, error } = await props.fetchItems()
      if (error != undefined) {
        items.value = []
        loadingFailed.value = true
      } else {
        items.value = data
      }
      loading.value = false
    }
  }

  onMounted(loadItems)


  const actions = {
    edit(item: ItemType) {
      return new Promise<ItemType>((resolve, reject) => {
        form.value = {
          mode: "Edit",
          editItem: item,
          dialog: true,
          onSuccess: resolve,
          onClose: reject
        }
      }).then(
        // Resolve
        (item) => {
          const index = items.value.findIndex(({ id }) => item.id === id)
          if (index < 0) {
            console.error('Failed to find edited item in currently loaded items', item)
            return
          }
          console.info('Edited item', item, ` at index ${index}`)
          items.value.splice(index, 1)
          items.value.unshift(item)
          triggerRef(items) // required to trigger recomputation of depending properties
          feedback.value.show(props.itemRepr ? `${props.itemRepr(item)} updated` : 'Item updated', 'success')
          emit('itemEdited', item, index)
          return { item, index }
        },
        // Reject
        () => {
          console.info('Item edition was cancelled')
          return
        }
      ).finally(() => {
        form.value.dialog = false
      })
    },
    create() {
      return new Promise<ItemType>((resolve, reject) => {
        form.value = {
          mode: 'Create',
          editItem: undefined,
          dialog: true,
          onSuccess: resolve,
          onClose: reject
        }
      }).then(
        // Resolve
        (item) => {
          console.info('Created item', item)
          items.value.unshift(item)
          triggerRef(items)
          feedback.value.show(props.itemRepr ? `${props.itemRepr(item)} registered` : 'Item registered', 'success')
          emit('itemCreated', item, 0)
          return { item, index: 0 }
        },
        // Reject
        () => {
          console.log('Item creation was cancelled')
          return undefined
        }
      ).finally(() => {
        form.value.dialog = false
      })
    },
    async delete(item: ItemType) {
      const message = props.itemRepr
        ? `Are you sure you want to delete ${props.itemRepr(item)} ?`
        : 'Are you sure you want to delete this item ?'
      return await askConfirm({
        title: "Confirm deletion",
        message,
        payload: item
      }).then(({ isCanceled, data }) => {
        if (isCanceled) {
          console.log("Item deletion canceled")
          return undefined
        }
        else if (data !== undefined)
          return executeDelete(data)
      })
    }
  }

  async function executeDelete(item: ItemType) {
    if (!item) {
      console.error('Item to delete is undefined. Aborting.')
      return
    }
    const index = items.value.findIndex(({ id }) => item.id === id)
    if (index === -1) console.error("Failed to find item index")
    console.log(`Deleting item at index ${index}`, item)
    if (props.delete == undefined) {
      items.value.splice(index, 1)
      triggerRef(items)
    } else {
      const { error } = await props.delete(item)
      if (error != undefined) {
        switch (error.status) {
          case HttpStatusCode.NotFound:
            feedback.value.show('Deletion failed: record not found.', 'error')
            break
          case HttpStatusCode.BadRequest:
            feedback.value.show(`Deletion was not allowed: ${error.detail}`, 'error')
            break
          case HttpStatusCode.Forbidden:
            feedback.value.show('You are not granted rights to modify this item.', 'error')
            break
          case HttpStatusCode.InternalServerError:
            feedback.value.show('An unexpected error occurred.', 'error')
        }
        return undefined
      } else {
        items.value.splice(index, 1)
        triggerRef(items);
        feedback.value.show('Item successfully deleted.', 'success')
        if (props.reloadOnDelete) {
          await loadItems()
        }
        return { item, index }
      }
    }
  }

  const feedback = ref<{
    model: boolean
    props: FeedbackProps
    show(text: string, color?: string): void
  }>({
    model: false,
    props: {
      text: '',
      color: undefined
    },
    show(text: string, color: string | undefined = undefined) {
      feedback.value.props = { text, color }
      feedback.value.model = true
    }
  })


  return {
    currentUser, feedback, actions, form, processedHeaders,
    loading, loadItems, loadingFailed,
  }
}