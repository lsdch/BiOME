import HomeView from '@/views/HomeView.vue';
import { RouteDefinition, RouterItem } from '.';
import routes from './routes';
import { useGuards } from './guards';

const { guardRole } = useGuards()

/** Route definitions meant to be displayed in navigation components */
export const navRoutes: RouterItem[] = [
  {
    label: "Home",
    path: '/',
    name: 'home',
    icon: "mdi-home",
    component: HomeView,
    meta: { subtitle: "Home" }
  },
  {
    icon: "mdi-folder-table",
    label: "Datasets",
    path: '/datasets',
    name: 'datasets',
    component: () => import('../views/datasets/DatasetsView.vue'),
    meta: { subtitle: "Datasets" }
  },
  {
    label: "Sampling",
    icon: "mdi-package-down",
    routes: [
      {
        label: "Sites",
        path: "/sites",
        name: "sites",
        icon: "mdi-map-marker-circle",
        component: () => import("@/views/location/SitesView.vue"),
        meta: {
          subtitle: "Sites",
          drawer: { temporary: true }
        }
      },
      {
        label: "Habitats",
        path: "/habitats",
        name: "habitats",
        icon: "mdi-image-filter-hdr-outline",
        component: () => import("@/views/location/HabitatsView.vue"),
        meta: { subtitle: "Habitats" }
      },
      {
        label: "Abiotic parameters",
        path: "/abiotic-parameters",
        name: "abiotic-parameters",
        icon: "mdi-gauge",
        component: () => import("@/views/sampling/AbioticParametersView.vue"),
        meta: { subtitle: "Abiotic parameters" }
      },
      {
        label: "Methods",
        path: "/sampling-methods",
        name: "sampling-methods",
        icon: "mdi-hook",
        component: () => import("@/views/sampling/SamplingMethodsView.vue")
      },
      {
        label: "Fixatives",
        path: "/fixatives",
        name: "fixatives",
        icon: "mdi-snowflake",
        component: () => import("@/views/sampling/FixativesView.vue")
      }
    ]
  },
  {
    label: "Samples",
    icon: "mdi-package-variant",
    routes: [
      {
        label: "Bio material",
        path: "/bio-material",
        name: "bio-material",
        icon: "mdi-package-variant",
        component: () => import("@/views/samples/BioMaterialView.vue")
      },
    ]
  },
  {
    label: "Sequences",
    icon: "mdi-dna",
    routes: [
      {
        label: "Genes",
        path: "/genes",
        name: "genes",
        icon: "mdi-tag",
        component: () => import("@/views/sequences/GenesView.vue")
      }
    ]
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
        component: () => import("../views/people/PersonView.vue"),
        meta: { subtitle: "Persons" }
      },
      {
        label: "Institutions",
        path: "/institutions",
        name: "institutions",
        icon: "mdi-domain",
        component: () => import("../views/people/InstitutionView.vue"),
        meta: { subtitle: "Institutions" }
      },
      {
        label: "Programs",
        path: "/programs",
        name: "programs",
        icon: "mdi-notebook",
        component: () => import("@/views/events/ProgramsView.vue"),
        meta: { subtitle: "Programs" }
      }
    ]
  },
  {
    icon: "mdi-family-tree",
    label: "Taxonomy",
    path: '/taxonomy',
    name: 'taxonomy',
    component: () => import('../views/taxonomy/TaxonomyView.vue'),
    beforeEnter: (to, from) => {
      if (from.path === to.path) {
        return false
      }
      return true
    },
    meta: { subtitle: "Taxonomy" }
  },
  {
    label: "Admin",
    icon: "mdi-cog",
    routes: [
      guardRole('Admin',
        {
          label: "Account requests",
          path: "/admin/account-requests",
          name: "account-requests",
          icon: "mdi-account-plus",
          component: () => import("@/views/registration/AccountsPendingView.vue"),
          meta: { subtitle: "Account requests" }
        }),
      routes.settings,
    ]
  },
]

export const navRouteDefinitions = navRoutes.reduce((acc, current) => {
  if (current.routes) {
    return acc.concat(current.routes)
  } else {
    acc.unshift(current as RouteDefinition)
    return acc
  }
}, new Array<RouteDefinition>)