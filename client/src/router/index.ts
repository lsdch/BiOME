import { InstanceSettings, UserRole } from '@/api'
import NotFound from '@/components/navigation/NotFound.vue'
import { nextTick } from "vue"
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useGuards } from './guards'

import routes, { accountRoutes } from './routes'
import { navRouteDefinitions } from './nav'

export * from './nav'


export type RouteNavDefinition = {
  label: string
  icon: string
  granted?: UserRole
}

export type RouteDefinition = RouteRecordRaw & RouteNavDefinition


type Route = RouteDefinition & { routes?: undefined }

type RouteGroup = Readonly<RouteNavDefinition & { routes: RouteDefinition[] }>

export type RouterItem = Route | RouteGroup

const { guardRole } = useGuards()


function setupRouter(settings: InstanceSettings) {
  function makeTitle(title: string) {
    return `${settings.name} | ${title}`
  }

  const router = createRouter({
    history: createWebHistory(),
    routes: [
      {
        path: '/docs/api',
        name: 'api-docs',
        component: () => import('../views/APIDocs.vue'),
        meta: {
          title: makeTitle("API docs")
        }
      },
      {
        path: "/init",
        name: "init",
        component: () => import("../views/InitialSetup.vue"),
        meta: { hideNavbar: true }
      },
      guardRole('Admin', {
        path: "/taxonomy/import",
        name: "import-GBIF",
        component: () => import("../views/taxonomy/GBIFImportView.vue")
      }),
      guardRole('Contributor', {
        path: "/location/import-sites",
        name: "import-sites",
        meta: {
          drawer: {
            temporary: true
          }
        },
        component: () => import("../views/location/SiteImportView.vue")
      }),
      {
        path: "/datasets/:slug",
        name: "dataset-item",
        component: () => import('@/views/datasets/DatasetItemView.vue')
      },
      {
        path: "/sites/:code",
        name: "site-item",
        component: () => import('@/views/location/SiteItemView.vue')
      },
      { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
      ...Object.values(routes),
      ...navRouteDefinitions
    ]
  })

  router.afterEach((to, _from) => {
    // Use next tick to handle router history correctly
    // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
    nextTick(() => {
      document.title = to.meta?.title ?? settings.name;
    });
  });

  return router
}

export default setupRouter
