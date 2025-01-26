import "@tanstack/vue-query"

import { Register } from "@tanstack/vue-query"
import { ErrorModel } from "@/api"

declare module "@tanstack/vue-query" {
  interface Register {
    defaultError: ErrorModel
  }
}