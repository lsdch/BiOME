import { ApiError } from "@/api"
import { ConfirmDialogKey } from "@/injection"
import { useUserStore } from "@/stores/user"
import { HttpStatusCode } from "axios"
import { ComputedRef, MaybeRef, ModelRef, computed, inject, onMounted, ref } from "vue"
import { FeedbackProps } from "../CRUDFeedback.vue"
import { Mode } from "../forms/form"


export type ToolbarProps = {
  /**
   * Table toolbar title
   */
  title: string
  /**
   * Table icon
   */
  icon?: string
  /**
   * Whether table search bar is displayed by default, or must be toggled
   */
  togglableSearch?: boolean
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
  toolbar: ToolbarProps
  /**
   * Display column with update/delete controls
   */
  showActions?: boolean
  /**
   * Short string representation of a table item to display in UI
   */
  itemRepr?: (item: ItemType) => string
  /**
   * API call to populate table items
   */
  fetchItems?: () => Promise<ItemType[]>
  /**
   * API call to delete an item
   */
  delete?: (item: ItemType) => Promise<ItemType>
  /**
   * Reload all items after deleting one
   */
  reloadOnDelete?: boolean
}



type FormSlotScope<ItemType extends { id: string }> = {
  dialog: boolean,
  mode: Mode,
  editItem?: MaybeRef<ItemType>,
  onSuccess: (_item: ItemType) => void,
  onClose: () => void,
}


export function useTable<ItemType extends { id: string }>(
  items: ModelRef<ItemType[]>,
  props: TableProps<ItemType>
) {

  const { user: currentUser } = useUserStore()

  const form = ref<FormSlotScope<ItemType>>({
    dialog: false,
    mode: 'Create',
    editItem: undefined,
    onSuccess: (_item: ItemType) => { },
    onClose: () => { }
  })

  const processedHeaders: ComputedRef<CRUDTableHeaders> = computed((): CRUDTableHeader[] => {
    return props.showActions && currentUser !== undefined && currentUser.role !== "Visitor"
      ? props.headers.concat([{ title: 'Actions', key: 'actions', sortable: false, align: 'end' }])
      : props.headers
  })

  const loading = ref(true)
  async function loadItems() {
    loading.value = true
    if (props.fetchItems) {
      items.value = await props.fetchItems()
    }
    loading.value = false
  }

  onMounted(loadItems)

  const confirmDelete = inject(ConfirmDialogKey)

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
          feedback.value.show(props.itemRepr ? `${props.itemRepr(item)} updated` : 'Item updated', 'success')
        },
        // Reject
        () => { console.info('Item edition was cancelled') }
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
        (item): void => {
          console.info('Created item', item)
          items.value.unshift(item)
          feedback.value.show(props.itemRepr ? `${props.itemRepr(item)} registered` : 'Item registered', 'success')
        },
        // Reject
        () => { console.info('Item creation was cancelled') }
      ).finally(() => {
        form.value.dialog = false
      })
    },
    async delete(item: ItemType) {
      const message = props.itemRepr
        ? `Are you sure you want to delete ${props.itemRepr(item)} ?`
        : 'Are you sure you want to delete this item ?'
      console.log(item)
      await confirmDelete?.({
        title: "Confirm deletion",
        message,
        data: item
      }).then(({ isCanceled, data }) => {
        console.log(data)
        if (isCanceled)
          console.info("Item deletion canceled")
        else if (data !== undefined)
          executeDelete(data)
      })
    }
  }

  async function executeDelete(item: ItemType) {
    if (!item) {
      console.error('Item to delete is undefined. Aborting.')
      return
    }
    console.log(`Deleting item`, item)
    const index = items.value.indexOf(item)

    if (props.delete == undefined) {
      items.value.splice(index, 1)
    } else {
      props.delete(item)
        .then(async () => {
          items.value.splice(index, 1)
          feedback.value.show('Item successfully deleted.', 'success')
          if (props.reloadOnDelete) {
            await loadItems()
          }
        })
        .catch((err: ApiError) => {
          switch (err.status) {
            case HttpStatusCode.NotFound:
              feedback.value.show('Deletion failed: record not found.', 'error')
              break
            case HttpStatusCode.BadRequest:
              feedback.value.show(`Deletion was not allowed: ${err.message}`, 'error')
              break
            case HttpStatusCode.Forbidden:
              feedback.value.show('You are not granted rights to modify this item.', 'error')
              break
            case HttpStatusCode.InternalServerError:
              feedback.value.show('An unexpected error occurred.', 'error')
          }
        })
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


  return { currentUser, feedback, actions, form, processedHeaders, loading, loadItems }
}