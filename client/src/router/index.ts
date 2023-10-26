import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'

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
      component: () => import('../views/SwaggerDocs.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/auth/LoginView.vue'),
      meta: { hideNavbar: true }
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('../views/auth/SignUpView.vue'),
      meta: { hideNavbar: true }
    },
    {
      path: '/users/password-reset',
      name: 'password-reset',
      component: () => import('../views/auth/PasswordResetView.vue'),
      meta: { hideNavbar: true }
    },
    ...routeGroups.reduce((acc, current) => acc.concat(current.routes), <RouteDefinition[]>[])
  ]
})

export default router
