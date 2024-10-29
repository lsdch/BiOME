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
    label: "Locations", icon: "mdi-map-marker-circle", routes: [
      {
        label: "Sites",
        path: "/sites",
        name: "sites",
        icon: "mdi-map-marker-radius",
        component: () => import("@/views/location/SitesView.vue"),
        meta: {
          subtitle: "Sites",
          drawer: { temporary: true }
        }
      },
      {
        label: "Habitats",
        path: "/location/habitats",
        name: "habitats",
        icon: "mdi-image-filter-hdr-outline",
        component: () => import("@/views/location/HabitatsView.vue"),
        meta: { subtitle: "Habitats" }
      },
    ]
  },
  {
    icon: "mdi-family-tree",
    label: "Taxonomy",
    path: '/taxonomy',
    name: 'taxonomy',
    component: () => import('../views/taxonomy/TaxonomyView.vue'),
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