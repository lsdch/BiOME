import { DeletePersonData, ErrorModel } from "@/api"
import { useAppConfirmDialog } from "@/composables/confirm_dialog"
import { useUserStore } from "@/stores/user"
import { OptionsLegacyParser, RequestResult } from "@hey-api/client-fetch"
import { StatusCodes } from "http-status-codes"
import { ComputedRef, MaybeRef, ModelRef, Ref, computed, onMounted, ref, triggerRef, watch } from "vue"
import { FeedbackProps } from "../CRUDFeedback.vue"
import { Mode } from "../forms/form"
import { listAbioticParametersQueryKey } from "@/api/gen/@tanstack/vue-query.gen"
import { QueryObserverResult, RefetchOptions, UndefinedInitialQueryOptions, useMutation, UseMutationOptions, useQuery } from "@tanstack/vue-query"
import { storeToRefs } from "pinia"



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


export type TableProps<ItemType extends {}, ItemsQueryData extends {}, ItemsDeleteData extends {}> = {
  /**
   * Entity name to display as title
   */
  entityName: string
  /**
   * Datatable headers definition
   */
  headers: CRUDTableHeader<ItemType>[]
  /**
   * Table toolbar configuration
   */
  toolbar?: ToolbarProps | false
  /**
   * Short string representation of a table item to display in UI
   */
  itemRepr?: (item: ItemType) => string
  /**
   * API call to populate table items
   */
  fetchItems?: (options?: OptionsLegacyParser<ItemsQueryData>) => UndefinedInitialQueryOptions<ItemType[], ErrorModel, ItemType[], any>
  /**
   * API call to delete an item
   */
  delete?: {
    mutation: (options?: ItemsDeleteData) => UseMutationOptions<
      ItemType,
      ErrorModel,
      ItemsDeleteData
    >,
    params: (item: ItemType) => ItemsDeleteData,
    fullReload?: boolean
  }
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


export function useTable<
  ItemType extends { id: string },
  ItemsQueryData extends {},
  ItemsDeleteData extends {}
>(
  items: ModelRef<ItemType[]>,
  props: TableProps<ItemType, ItemsQueryData, ItemsDeleteData>,
  emit: TableEmits<ItemType>
) {

  const { user: currentUser } = storeToRefs(useUserStore())
  const { askConfirm } = useAppConfirmDialog()

  const form = ref<FormSlotScope<ItemType>>({
    dialog: false,
    mode: 'Create',
    editItem: undefined,
    onSuccess: (_item: ItemType) => { },
    onClose: () => { }
  })

  const processedHeaders = computed((): CRUDTableHeader<ItemType>[] => {
    return props.headers.filter(({ hide }) => {
      return !hide?.value
    })
  }) as ComputedRef<DataTableHeader[]>


  // Items fetching
  const { data, error, isFetching, isSuccess, refetch } = props.fetchItems
    ? useQuery({ ...props.fetchItems(), enabled: () => !items.value.length, initialData: [] })
    : {
      data: ref<ItemType[]>(items.value ?? []) as Ref<ItemType[]>,
      error: ref<ErrorModel>(),
      isFetching: ref(false),
      isSuccess: ref(true),
      refetch: () => new Promise<void>((resolve, reject) => {
        return resolve()
      })
    }

  watch(data, (data) => items.value = [...data], { immediate: true })

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
          items.value![index] = item
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
          // items.value = [item, ...items.value]
          items.value!.unshift(item)
          triggerRef(items) // required to trigger recomputation of depending properties
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
      }).then(async ({ isCanceled, data }) => {
        if (isCanceled) {
          console.log("Item deletion canceled")
          return undefined
        }

        const index = items.value!.findIndex(({ id }) => item.id === id)
        if (index === -1) console.error("Failed to find item index")
        console.log(`Deleting item at index ${index}`, item)

        if (props.delete == undefined) {
          items.value!.splice(index, 1)
          triggerRef(items)
        } else {
          await mutateAsync(props.delete.params(item))
          if (deleteSuccess.value) {
            items.value!.splice(index, 1)
            triggerRef(items);
          } else {
            return undefined
          }
        }
        return { item, index }
      })
    }
  }

  // Delete mutation
  const { mutateAsync, isSuccess: deleteSuccess } = useMutation({
    ...props.delete?.mutation(),
    onSuccess(deleted) {
      feedback.value.show('Item successfully deleted.', 'success')
      if (props.delete?.fullReload) refetch()
    },
    onError(error) {
      switch (error.status) {
        case StatusCodes.NOT_FOUND:
          feedback.value.show('Deletion failed: record not found.', 'error')
          break
        case StatusCodes.BAD_REQUEST:
          feedback.value.show(`Deletion was not allowed: ${error.detail}`, 'error')
          break
        case StatusCodes.FORBIDDEN:
          feedback.value.show('You are not granted rights to modify this item.', 'error')
          break
        case StatusCodes.INTERNAL_SERVER_ERROR:
          feedback.value.show('An unexpected error occurred.', 'error')
      }
    }
  })

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
    loading: isFetching,
    loadItems: refetch,
    error,
  }
}