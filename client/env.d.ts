/// <reference types="vite/client" />
import "vue-router"

declare module 'swagger-ui';

interface ImportMetaEnv {
  readonly VITE_APP_NAME: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
export { }
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    hideNavbar?: boolean
    drawer?: {
      temporary?: boolean
    }
  }
}

