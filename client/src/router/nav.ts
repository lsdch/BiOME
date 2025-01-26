import HomeView from '@/views/HomeView.vue';
import { Divider, RouteDefinition, RouterItem } from '.';
import routes from './routes';
import { useGuards } from './guards';

const { guardRole } = useGuards()



export function isDivider(item: RouterItem | Divider): item is Divider {
  return item === "divider"
}

/** Route definitions meant to be displayed in navigation components */
export const navRoutes: (RouterItem | Divider)[] = [
  {
    label: "Home",
    path: '/',
    name: 'home',
    icon: "mdi-home",
    component: HomeView,
    meta: { subtitle: "Home" }
  },
  {
    label: "Datasets",
    icon: "mdi-folder-table",
    routes: [
      {
        label: "Sites",
        path: '/datasets/sites',
        name: 'site-datasets',
        icon: 'mdi-map-marker-circle',
        component: () => import('../views/datasets/SiteDatasetsView.vue'),
        meta: { subtitle: "Site datasets" }
      },
      {
        label: "Occurrences",
        path: '/datasets/occurrences',
        name: 'occurrence-datasets',
        icon: 'mdi-crosshairs-gps',
        component: () => import('../views/datasets/OccurrenceDatasetsView.vue'),
        meta: { subtitle: "Occurrence datasets" }
      },
      {
        label: "Sequences",
        path: '/datasets/sequences',
        name: 'seq-datasets',
        icon: 'mdi-dna',
        component: () => import('../views/datasets/SeqDatasetsView.vue'),
        meta: { subtitle: "Sequence datasets" }
      }
    ]
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
        meta: { subtitle: "Habitats" },
        props: {
          density: "compact"
        }
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
        label: "Sequences",
        path: "/sequences",
        name: "sequences",
        icon: "mdi-dna",
        component: () => import("@/views/sequences/SequencesView.vue")
      },
      {
        label: "Genes",
        path: "/genes",
        name: "genes",
        icon: "mdi-tag",
        component: () => import("@/views/sequences/GenesView.vue")
      },
      {
        label: "Seq. databases",
        path: "/seq-databases",
        name: "seq-databases",
        icon: "mdi-database-sync",
        component: () => import("@/views/sequences/SeqDatabasesView.vue")
      }
    ]
  },
  {
    label: "DNA sequencing",
    icon: "mdi-flask",
    routes: []
  },
  "divider",
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
    label: "Bibliography",
    icon: 'mdi-newspaper-variant-multiple',
    name: "bibliography",
    path: '/articles',
    component: () => import('@/views/references/ArticlesView.vue')
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
  if (isDivider(current)) {
    return acc
  }
  if (current.routes) {
    return acc.concat(current.routes)
  } else {
    acc.unshift(current as RouteDefinition)
    return acc
  }
}, new Array<RouteDefinition>)