import { useConfirmDialog } from "@vueuse/core"
import { ref } from "vue"


export type ConfirmDialogContent<P> = {
  title: string
  message?: string
  payload?: P
}

const confirmDialog = ref<ConfirmDialogContent<any>>({
  title: '',
  message: undefined,
  payload: undefined
})

const { isRevealed, confirm, cancel, reveal, onReveal } = useConfirmDialog<ConfirmDialogContent<any>>()

onReveal((data) => confirmDialog.value = data)

export function useAppConfirmDialog() {
  function askConfirm<P>(data: ConfirmDialogContent<P>) {
    return reveal(data)
  }

  return { isRevealed, confirm: () => confirm(confirmDialog.value.payload), cancel, askConfirm, content: confirmDialog }
}
