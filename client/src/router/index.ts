import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import NotFound from '@/components/NotFound.vue'
import { nextTick } from "vue"


// Default app title to display
const DEFAULT_TITLE = import.meta.env.VITE_APP_NAME

function makeTitle(title: string) {
  return `${DEFAULT_TITLE} - ${title}`
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
    component: () => import('../views/taxonomy/TaxonomyMain.vue')
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
  }
]



const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/docs/api',
      name: 'api-docs',
      component: () => import('../views/SwaggerDocs.vue'),
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
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
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
    document.title = to.meta?.title ?? DEFAULT_TITLE;
  });
});

export default router
