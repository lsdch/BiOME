import { UseConfirmDialogRevealResult } from "@vueuse/core"
import { InjectionKey } from "vue"
import { ConfirmDialogProps } from "./components/toolkit/ui/ConfirmDialog.vue"

export const ConfirmDialogKey = Symbol() as InjectionKey<<T>(dialog: ConfirmDialogProps<T>) => Promise<UseConfirmDialogRevealResult<T, undefined>>>