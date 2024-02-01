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

type RouteDefinition = RouteRecordRaw & {
  label: string
  // icon?: string
}

type RouteGroup = {
  readonly name?: string // unused if only one route
  readonly icon: string
  readonly routes: RouteDefinition[]
}

/** Route definitions meant to be displayed in navigation components */
export const routeGroups: RouteGroup[] = [
  {
    icon: "mdi-home",
    routes: [{
      label: "Home",
      path: '/',
      name: 'home',
      component: HomeView
    }]
  },
  {
    name: "Taxonomy",
    icon: "mdi-graph",
    routes: [
      {
        label: "Taxonomy",
        path: '/taxonomy',
        name: 'taxonomy',
        component: () => import('../views/Taxonomy/TaxonomyMain.vue')
      }
    ]
  },
  {
    name: "Tmp",
    icon: "mdi-flask",
    routes: [
      {
        label: "About",
        path: '/about',
        name: 'about',
        component: () => import('../views/AboutView.vue')
      },
      {
        label: "Other",
        path: '/about',
        name: 'about',
        component: () => import('../views/AboutView.vue')
      },
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
    ...routeGroups.reduce((acc, current) => acc.concat(current.routes), <RouteDefinition[]>[])
  ]
})

router.afterEach((to, from) => {
  // Use next tick to handle router history correctly
  // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
  nextTick(() => {
    document.title = to.meta?.title ?? DEFAULT_TITLE;
  });
});

export default router
