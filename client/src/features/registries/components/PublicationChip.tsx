import { Publication } from '@/api'
import { VChip } from 'vuetify/components'

export type PublicationChipProps = {
  publication: Publication
} & VChip['$props']

export function PublicationChip({ publication, ...chipProps }: PublicationChipProps) {
  return (
    <v-menu
      location="top start"
      origin="top start"
      transition="scale-transition"
      open-on-focus={false}
      open-on-click={true}
    >
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip {...{ ...props, ...chipProps }} style={'max-width: 300px'} class="text-truncate">
            <span class="text-truncate" style={{
              'min-width': '0px',
              'max-width': '100%',
            }}>{Publication.toString(publication)}</span>
            </v-chip>
        ),
        default: () => (
          <v-card
            title={publication.title ?? publication.verbatim ?? 'Untitled publication'}
            subtitle={publication.journal ?? 'Unknown journal'}
            class="small-card-title bg-surface-light"
            density="compact"
            max-width={600}
          >
            {{
              append: () => <v-chip label text={publication.year?.toString()} />,
              default: () => <v-card-text>{publication.authors?.join(', ')}</v-card-text>,
              actions: () =>
                publication.doi ? (
                  <v-card-actions>
                    {publication.doi ? (
                      <a href={Publication.linkDOI(publication)}>{Publication.linkDOI(publication)}</a>
                    ) : null}
                  </v-card-actions>
                ) : null
            }}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default PublicationChip
