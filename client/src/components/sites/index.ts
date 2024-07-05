import { LocationService } from "@/api";
import { handleErrors } from "@/api/responses";
import { onKeyDown, onKeyStroke, useDebouncedRefHistory, useKeyModifier, useMagicKeys, whenever } from "@vueuse/core";
import { parse } from "papaparse";
import { computed, Ref, ref } from "vue";
import { Schema, useSchema } from "../toolkit/forms/schema";
import { Errors, indexErrors } from "../toolkit/validation";

export async function useAccessPoints() {
  const accessPoints = ref(await fetch())

  function fetch() {
    return LocationService.getAccessPoints().then(handleErrors(err => console.error("Failed to fetch access points", err)))
  }

  async function refresh() {
    accessPoints.value = await fetch()
  }

  return { refresh, accessPoints }
}


export type Selection = {
  start: {
    row: number
    col: number
  }
  end?: {
    row: number
    col: number
  }
}


function* range(start: number, end: number, step = 1) {
  let i = start
  while (i < end) {
    yield i;
    i += step
  }
}

export function useSpreadsheet<Item extends {} & { errors?: Errors<string> }>(
  setValue: (row: number, col: number, value: any) => void,
  schema?: Schema
) {

  const items: Ref<Partial<Item>[]> = ref(Array.from(Array(10), () => ({})))

  const { validateAll } = schema ? useSchema(schema) : { validateAll: undefined }

  const selection = ref<Selection>()
  const editing = ref<{ row: number; col: number }>()
  const editingItem = computed(() => editing.value ? items.value[editing.value.row] : undefined)
  const dragging = ref(false)

  const { undo, redo } = useDebouncedRefHistory(items, { deep: true, debounce: 500 })

  const magicKeys = useMagicKeys()

  const CtrlZ = magicKeys['Ctrl+Z']
  const CtrlY = magicKeys['Ctrl+Y']

  whenever(CtrlZ, undo)
  whenever(CtrlY, redo)

  onKeyStroke('ArrowDown', () => moveSelection({ row: 1 }))
  onKeyStroke('ArrowUp', () => moveSelection({ row: -1 }))
  onKeyStroke('ArrowRight', () => moveSelection({ col: 1 }))
  onKeyStroke('ArrowLeft', () => moveSelection({ col: -1 }))
  onKeyStroke('Delete', () => {
    if (selection.value === undefined) return
    const { start, end } = selection.value
    const e = end ?? start
    const cols = [start.col, e.col].sort()
    const rows = [start.row, e.row].sort()
    Array.from(range(rows[0], rows[1] + 1)).forEach(
      (row) => {
        Array.from(range(cols[0], cols[1] + 1)).forEach(
          (col) => {
            setValue(row, col, undefined)
          }
        )
        onEdited(items.value[row])
      }
    )
  })

  const shiftMod = useKeyModifier('Shift')
  onKeyDown('Tab', (e) => {
    if (selection.value === undefined) return
    e.preventDefault()
    if (editing.value) onEdited(items.value[editing.value.row])

    shiftMod.value
      ? moveSelection({ col: - 1, row: selection.value?.start.col === 0 ? -1 : 0 })
      : moveSelection({
        col: 1,
        row: selection.value.start.col == nCol.value - 1 ? 1 : 0
      })
    selection.value.end = undefined
  }
  )

  onKeyStroke('Enter', (e) => {
    editingItem.value === undefined
      ? (editing.value = selection.value?.start) && e.stopPropagation()
      : onEdited(editingItem.value)
  })

  onKeyStroke('Escape', () => {
    editingItem.value
      ? onEdited(editingItem.value)
      : selection.value = undefined
  })

  function moveSelection({ row, col }: { row?: number; col?: number }) {
    if (selection.value === undefined || editing.value !== undefined) return
    const cell = selection.value.start
    cell.row = Math.min(Math.max(cell.row + (row ?? 0), 0), items.value.length - 1)

    cell.col += col ?? 0
    cell.col = cell.col > 2 ? 0 : cell.col < 0 ? 2 : cell.col
  }

  function select(row: number, col: number) {
    dragging.value = true
    selection.value = { start: { row, col } }
  }
  function selectEnd(row: number, col: number) {
    if (selection.value === undefined) return
    selection.value.end = { row, col }
  }

  function isInRange(x: number, range: [number, number]) {
    range.sort()
    return x >= range[0] && x <= range[1]
  }

  function isSelected(row: number, col?: number) {
    if (selection.value === undefined) return false
    return (
      isInRange(row, [
        selection.value.start.row,
        selection.value.end?.row ?? selection.value.start.row
      ]) && (
        col === undefined ||
        isInRange(col, [
          selection.value.start.col,
          selection.value.end?.col ?? selection.value.start.col
        ])
      )
    )
  }


  function isEditing(row: number, col: number) {
    return editing.value?.row == row && editing.value.col == col
  }

  function handlePaste(event: ClipboardEvent) {
    const pasteData = event.clipboardData?.items
    if (pasteData === undefined || selection.value === undefined) return
    const sel = selection.value
    const data = Array.from(pasteData)?.find((item) => item.type === 'text/plain')
    data?.getAsString((s) => {
      const { data } = parse<Array<string | number>>(s, {
        header: false,
        delimiter: '\t',
        dynamicTyping: true
      })
      const startRow = Math.min(sel.start.row, sel.end?.row ?? sel.start.row)
      const startCol = Math.min(sel.start.col, sel.end?.col ?? sel.start.col)
      const needRows = startRow + data.length
      if (items.value.length < needRows) {
        const newRows = ref(Array.from(Array(needRows - items.value.length), () => ({})))
        items.value = items.value.concat(newRows.value)
      }
      data.forEach((arr, index) => {
        if (arr.length == 1 && arr[0] === null) return
        arr.forEach((v, colOffset) => {
          setValue(startRow + index, startCol + colOffset, v)
        })
        if (validateAll) {
          const item = items.value[startRow + index]
          item.errors = indexErrors(validateAll(items.value[startRow + index]))
        }
      })
    })
  }

  /** Number of editable columns */
  const nCol = ref(0)
  /** Column ID generator */
  function colGen() {
    return (nCol.value += 1) - 1
  }

  /**
   * Marks a column as a spreadsheet column,
   * making its cells selectable and editable
   */
  function cellHeader(header: DataTableHeader & {
    class?: string,
    key: keyof Omit<Item, 'errors'>
  }): DataTableHeader {
    const col = colGen()
    return {
      ...header,
      cellProps: ({ index, item }: { index: number, item: Item }) => ({
        onPaste: handlePaste,
        class: [
          'cell',
          header.class,
          {
            ['error-cell']: item.errors?.[header.key] !== undefined,
            editing: isEditing(index, col),
            selected: isSelected(index, col),
            right:
              selection.value &&
              col === Math.max(selection.value.start.col, selection.value.end?.col ?? 0),
            left:
              selection.value &&
              col ===
              Math.min(
                selection.value.start.col,
                selection.value.end?.col ?? selection.value.start.col
              ),
            top:
              selection.value &&
              index ===
              Math.min(
                selection.value.start.row,
                selection.value.end?.row ?? selection.value.start.row
              ),
            bottom:
              selection.value &&
              index === Math.max(selection.value.start.row, selection.value.end?.row ?? 0)
          }
        ]
      })
    }
  }

  function edit(row: number, col: number) {
    editing.value = { row, col }
  }

  function isEmpty(item: Partial<Item>) {
    return Object.entries(item)
      .filter(([k, v]) => k !== 'errors' && v !== undefined)
      .length === 0
  }

  function onEdited(item: Partial<Item>) {
    editing.value = undefined
    if (validateAll !== undefined) {
      item.errors = !isEmpty(item) ? indexErrors(validateAll(item)) : undefined
    }
  }



  return { items, selection, editing, dragging, edit, onEdited, cellHeader, handlePaste, isSelected, isEditing, moveSelection, select, selectEnd, isEmpty }
}