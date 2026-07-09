import HomeView from '@/views/HomeView.vue'
import { Divider, RouteDefinition, RouterItem } from '.'
import { useGuards } from './guards'
import routes from './routes'
import { InstanceSettings } from '@/api'

const { guardRole } = useGuards()

export function isDivider(item: RouterItem | Divider): item is Divider {
  return item === 'divider'
}

/** Route definitions meant to be displayed in navigation components */
export function navRoutes(settings: InstanceSettings): (RouterItem | Divider)[] {
  return [
    {
      label: 'Home',
      path: '/',
      name: 'home',
      icon: 'mdi-home',
      component: HomeView,
      meta: { title: 'Home' }
    },
    {
      label: 'Mapping tool',
      path: '/mapping',
      name: 'mapping',
      icon: 'mdi-map-marker-circle',
      component: () => import('@/features/cartography/views/MappingToolView.vue'),
      meta: {
        title: 'Mapping tool',
        drawer: { temporary: true }
      }
    },
    {
      label: 'Datasets',
      icon: 'mdi-folder-table',
      routes: [
        // {
        //   label: "Sites",
        //   path: '/datasets/sites',
        //   name: 'site-datasets',
        //   icon: 'mdi-map-marker-circle',
        //   component: () => import('@/features/datasets/views/SiteDatasetsView.vue'),
        //   meta: { title: "Site datasets" }
        // },
        {
          label: 'Occurrences',
          path: '/datasets/occurrences',
          name: 'occurrence-datasets',
          icon: 'mdi-crosshairs-gps',
          component: () => import('@/features/datasets/views/OccurrenceDatasetsView.vue'),
          meta: { title: 'Occurrence datasets' }
        },
        {
          label: 'Sequences',
          path: '/datasets/sequences',
          name: 'seq-datasets',
          icon: 'mdi-dna',
          component: () => import('@/features/datasets/views/SeqDatasetsView.vue'),
          meta: { title: 'Sequence datasets' }
        }
        // {
        //   label: "Research programs",
        //   path: "/programs",
        //   name: "programs",
        //   icon: "mdi-notebook",
        //   component: () => import("@/features/datasets/views/ProgramsView.vue"),
        //   meta: { title: "Programs" }
        // }
      ]
    },
    {
      label: 'Occurrences',
      path: '/occurrences',
      name: 'occurrences',
      icon: 'mdi-package-variant',
      component: () => import('@/features/occurrences/views/OccurrencesTableView.vue'),
      meta: { title: 'Occurrences' }
    },
    {
      label: 'Sequences',
      path: '/sequences',
      name: 'sequences',
      icon: 'mdi-dna',
      hidden: !settings.molecular_data_enabled,
      component: () => import('@/features/sequences/views/SequencesView.vue'),
      meta: { title: 'Sequences' }
    },
    guardRole('Admin', {
      label: 'Data inputs',
      path: '/import',
      name: 'import',
      icon: 'mdi-file-upload',
      component: () => import('@/features/import/views/DataImportView.vue'),
      meta: { title: 'Data import' }
    }),
    {
      label: 'Data inputs',
      icon: 'mdi-file-upload',
      granted: 'Contributor',
      routes: [
        guardRole('Contributor', {
          label: 'Register item',
          path: '/import/item',
          name: 'import-item',
          icon: 'mdi-package-variant-plus',
          component: () => import('@/features/import/views/CreateOccurrencesView.vue'),
          meta: { title: 'Create sampling/occurrences' }
        }),
        guardRole('Contributor', {
          label: 'Import batch',
          path: '/import/batch',
          name: 'import-batch',
          icon: 'mdi-file-upload',
          component: () => import('@/features/import/views/ImportBatchView.vue'),
          meta: { title: 'Import batch' }
        })
      ]
    },
    'divider',
    // {
    //   label: "DNA sequencing",
    //   icon: "mdi-flask",
    //   routes: []
    // },
    {
      icon: 'mdi-family-tree',
      label: 'Taxonomy',
      path: '/taxonomy',
      name: 'taxonomy',
      component: () => import('@/features/taxonomy/views/TaxonomyView.vue'),
      beforeEnter: (to, from) => {
        if (from.path === to.path) {
          return false
        }
        return true
      },
      meta: { title: 'Taxonomy' }
    },
    {
      label: 'Users',
      path: '/users',
      name: 'users',
      icon: 'mdi-account',
      component: () => import('@/features/people/views/UsersView.vue'),
      meta: { title: 'Persons' }
    },
    {
      label: 'Metadata registries',
      icon: 'mdi-book-alphabet',
      routes: [
        { subgroup: 'Sampling' },
        {
          label: 'Habitats',
          path: '/habitats',
          name: 'habitats',
          icon: 'mdi-image-filter-hdr-outline',
          component: () => import('@/features/registries/views/HabitatsView.vue'),
          meta: { title: 'Habitats' }
        },
        {
          label: 'Methods',
          path: '/sampling-methods',
          name: 'sampling-methods',
          icon: 'mdi-hook',
          component: () => import('@/features/registries/views/SamplingMethodsView.vue'),
          meta: { title: 'Sampling methods' }
        },
        {
          label: 'Fixatives',
          path: '/fixatives',
          name: 'fixatives',
          icon: 'mdi-snowflake',
          component: () => import('@/features/registries/views/FixativesView.vue'),
          meta: { title: 'Fixatives' }
        },
        {
          label: 'Abiotic parameters',
          path: '/abiotic-parameters',
          name: 'abiotic-parameters',
          icon: 'mdi-gauge',
          component: () => import('@/features/registries/views/AbioticParametersView.vue'),
          meta: { title: 'Abiotic parameters' }
        },
        { subgroup: 'Sequences', hidden: !settings.molecular_data_enabled },
        {
          label: 'Genes',
          path: '/genes',
          name: 'genes',
          icon: 'mdi-tag',
          hidden: !settings.molecular_data_enabled,
          component: () => import('@/features/sequences/views/GenesView.vue'),
          meta: { title: 'Genes' }
        },
        { subgroup: 'Sources' },
        {
          label: 'Bibliography',
          icon: 'mdi-newspaper-variant-multiple',
          name: 'bibliography',
          path: '/articles',
          component: () => import('@/features/registries/views/ArticlesView.vue'),
          meta: { title: 'Bibliography' }
        },
        {
          label: 'Data sources',
          path: '/data-sources',
          name: 'data-sources',
          icon: 'mdi-database-sync',
          component: () => import('@/features/registries/views/DataSourcesView.vue'),
          meta: { title: 'Data sources' }
        }
        // {
        //   label: "Collections",
        //   path: "/collections",
        //   name: "collections",
        //   icon: "mdi-library-shelves",
        //   component: () => import("@/features/registries/views/CollectionsView.vue")
        // }
      ]
    },
    {
      label: 'Admin',
      icon: 'mdi-cog',
      granted: 'Admin',
      routes: [
        guardRole('Admin', {
          label: 'Account requests',
          path: '/admin/account-requests',
          name: 'account-requests',
          icon: 'mdi-account-plus',
          component: () => import('@/views/accounts/AccountsPendingView.vue'),
          meta: { title: 'Account requests' }
        }),
        routes.settings
      ]
    }
  ]
}

export const navRouteDefinitions = (settings: InstanceSettings) =>
  navRoutes(settings).reduce((acc, current) => {
    if (isDivider(current) || current.hidden) {
      return acc
    }
    if (current.routes) {
      return acc.concat(current.routes.filter((r): r is RouteDefinition => !('subgroup' in r)))
    } else {
      acc.unshift(current as RouteDefinition)
      return acc
    }
  }, new Array<RouteDefinition>())
