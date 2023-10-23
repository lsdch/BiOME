/// <reference types="vite/client" />

declare module 'swagger-ui';

interface ImportMetaEnv {
  readonly VITE_APP_NAME: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}