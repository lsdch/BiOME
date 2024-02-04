import { CancelablePromise } from "@/api"
import { Component } from "vue"


export type ToolbarProps<ItemType> = {
  title: string
  icon?: string
  entityName: string
  itemRepr?: (item: ItemType) => string
  togglableSearch?: boolean
  form?: Component
}

export type TableProps<ItemType, ItemInputType, FetchList> = {
  headers: ReadonlyHeaders
  toolbarProps: ToolbarProps<ItemType>
  showActions?: boolean
  crud: {
    list: FetchList
    delete: (item: ItemType) => CancelablePromise<any>
    create?: (item: ItemInputType) => CancelablePromise<any>
    update: (item: ItemType) => CancelablePromise<any>
  }
}
