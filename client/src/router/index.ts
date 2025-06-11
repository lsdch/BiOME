import { InstanceSettings, UserRole } from '@/api'
import NotFound from '@/components/navigation/NotFound.vue'
import { nextTick } from "vue"
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { useGuards } from './guards'

import { getSiteDatasetOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { ComponentProps } from 'vue-component-type-helpers'
import { VListGroup, VListItem } from 'vuetify/components'
import { navRouteDefinitions } from './nav'
import routes from './routes'

export * from './nav'


export type RouteNavDefinition = {
  label: string
  icon: string
  granted?: UserRole
  itemProps?: ComponentProps<typeof VListItem>
}

export type Divider = "divider"
export type RouteDefinition = RouteRecordRaw & RouteNavDefinition
export type RouteSubgroup = { subgroup: string }
export type Route = RouteDefinition & { routes?: undefined }
export type RouteGroup = Readonly<RouteNavDefinition & { groupProps?: ComponentProps<typeof VListGroup>, routes: (RouteDefinition | RouteSubgroup)[] }>
export type RouterItem = Route | RouteGroup

const { guardRole } = useGuards()


function setupRouter(settings: InstanceSettings) {
  function makeTitle(subtitle?: string) {
    return subtitle ? `${settings.name} | ${subtitle}` : settings.name
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
        path: "/import/dataset",
        name: "import-dataset",
        meta: {
          drawer: {
            temporary: true
          }
        },
        component: () => import("../views/location/SiteImportView.vue")
      }),
      {
        path: "/datasets/sites/:slug",
        name: "site-dataset-item",
        component: () => import('@/views/datasets/SiteDatasetItemView.vue'),
        props: route => ({ slug: route.params.slug, query: getSiteDatasetOptions }),
      },
      {
        path: "/datasets/occurrences/:slug",
        name: "occurrence-dataset-item",
        component: () => import('@/views/datasets/OccurrenceDatasetItemView.vue'),
        props: route => ({ slug: route.params.slug }),
      },
      {
        path: "/sites/:code",
        name: "site-item",
        component: () => import('@/views/location/SiteItemView.vue'),
        props: true
      },
      {
        path: "/bio-material/:code",
        name: "biomat-item",
        component: () => import('@/views/samples/BiomatItemView.vue'),
        props: true
      },
      {
        path: "/sequences/:code",
        name: "sequence",
        component: () => import('@/views/sequences/SeqItemView.vue'),
        props: true
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
      document.title = to.meta?.title ?? makeTitle(to.meta.subtitle);
    });
  });

  return router
}

export default setupRouter
