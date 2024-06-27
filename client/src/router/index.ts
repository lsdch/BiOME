import NotFound from '@/components/navigation/NotFound.vue'
import { nextTick } from "vue"
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useGuards } from './guards'
import { useInstanceSettings } from '@/components/settings'


// Default app title to display
const { settings } = useInstanceSettings()
function makeTitle(title: string) {
  return `${settings.name} | ${title}`
}

export type RouteDefinition = RouteRecordRaw & {
  label: string
  icon?: string
}

type Route = RouteDefinition & { icon: string, routes?: undefined }

type RouteGroup = {
  readonly label: string
  readonly icon: string
  readonly routes: RouteDefinition[]
}

export type RouterItem = Route | RouteGroup



/** Route definitions meant to be displayed in navigation components */
export const routeGroups: RouterItem[] = [
  {
    label: "Home",
    path: '/',
    name: 'home',
    icon: "mdi-home",
    component: HomeView
  },
  {
    icon: "mdi-graph",
    label: "Taxonomy",
    path: '/taxonomy',
    name: 'taxonomy',
    component: () => import('../views/taxonomy/TaxonomyView.vue')
  },
  {
    label: "People",
    icon: "mdi-account-group",
    routes: [
      {
        label: "Persons",
        path: "/people",
        name: "people",
        icon: "mdi-account",
        component: () => import("../views/people/PersonView.vue")
      },
      {
        label: "Institutions",
        path: "/people/institutions",
        name: "institutions",
        icon: "mdi-domain",
        component: () => import("../views/people/InstitutionView.vue")
      }
    ]
  },
  {
    label: "Locations", icon: "mdi-map-marker-circle", routes: [
      {
        label: "Sites",
        path: "/location/sites",
        name: "sites",
        icon: "mdi-map-marker-radius",
        component: () => import("@/views/location/SitesView.vue"),
        meta: {
          drawer: {
            temporary: true
          }
        }
      },
      {
        label: "Habitats",
        path: "/location/habitats",
        name: "habitats",
        icon: "mdi-image-filter-hdr-outline",
        component: () => import("@/views/location/HabitatsView.vue")
      },
    ]
  }
]



function router() {
  const { guardRole } = useGuards()
  const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
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
        path: '/login',
        name: 'login',
        component: () => import('../views/auth/LoginView.vue'),
        // meta: { hideNavbar: true }
      },
      {
        path: '/signup',
        name: 'signup',
        component: () => import('../views/auth/SignUpView.vue'),
        // meta: { hideNavbar: true }
      },
      {
        path: '/password-reset/:token',
        name: 'password-reset',
        component: () => import('../views/auth/PasswordResetView.vue'),
        meta: { hideNavbar: true }
      },
      {
        path: '/email-confirmation',
        name: 'email-confirmation',
        component: () => import('../views/auth/EmailConfirmation.vue'),
        meta: { hideNavbar: true }
      },
      {
        path: "/init",
        name: "init",
        component: () => import("../views/InitialSetup.vue"),
        meta: { hideNavbar: true }
      },
      {
        path: "/account",
        name: "account",
        component: () => import("../views/AccountView.vue")
      },
      guardRole('Contributor', {
        path: "/location/import-sites",
        name: "import-sites",
        component: () => import("../views/location/SiteImportView.vue")
      }),
      { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
      guardRole('Admin', {
        path: '/settings',
        name: "app-settings",
        component: () => import("@/views/settings/AdminSettings.vue")
      }),
      ...routeGroups.reduce((acc, current) => {
        if (current.routes) {
          return acc.concat(current.routes)
        } else {
          acc.unshift(current as RouteDefinition)
          return acc
        }
      },
        new Array<RouteDefinition>)

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
export default router
