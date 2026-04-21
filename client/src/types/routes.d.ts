import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    /**
     * Route title to display in browser title bar
     */
    title?: string
    /**
     * Hide top navbar when this route is active
     */
    hideNavbar?: boolean
    /**
     * Navigation drawer configuration
     */
    drawer?: {
      /**
       * Drawer is hidden when not focused
       */
      temporary?: boolean
    }
  }
}
