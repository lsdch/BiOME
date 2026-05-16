import { InstanceSettings, UserRole } from '@/api'
import NotFound from '@/components/navigation/NotFound.vue'
import { nextTick } from 'vue'
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

export type Divider = 'divider'
export type RouteDefinition = RouteRecordRaw & RouteNavDefinition
export type RouteSubgroup = { subgroup: string }
export type Route = RouteDefinition & { routes?: undefined }
export type RouteGroup = Readonly<
  RouteNavDefinition & {
    groupProps?: ComponentProps<typeof VListGroup>
    routes: (RouteDefinition | RouteSubgroup)[]
  }
>
export type RouterItem = Route | RouteGroup

const { guardRole } = useGuards()

function setupRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      {
        path: '/docs/api',
        name: 'api-docs',
        component: () => import('@/views/APIDocs.vue'),
        meta: {
          title: 'API docs'
        }
      },
      {
        path: '/datasets/occurrences/:slug',
        name: 'occurrence-dataset-item',
        component: () =>
          import('@/features/datasets/views/occurrence/OccurrenceDatasetItemView.vue'),
        props: (route) => ({ slug: route.params.slug })
      },
      {
        path: '/sites/:code',
        name: 'site-item',
        component: () => import('@/features/site/views/SiteItemView.vue'),
        props: true
      },
      {
        path: '/occurrences/:code',
        name: 'occurrence-item',
        component: () => import('@/features/occurrences/views/OccurrenceItemView.vue'),
        props: true
      },
      {
        path: '/sequences/:code',
        name: 'sequence',
        component: () => import('@/features/sequences/views/SeqItemView.vue'),
        props: true
      },
      { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
      ...Object.values(routes),
      ...navRouteDefinitions
    ]
  })

  return router
}

export default setupRouter
