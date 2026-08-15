import { Article } from '@/api'
import { VChip } from 'vuetify/components'

export type ArticleChipProps = {
  article: Article
} & VChip['$props']

export function ArticleChip({ article, ...chipProps }: ArticleChipProps) {
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
            }}>{Article.toString(article)}</span>
            </v-chip>
        ),
        default: () => (
          <v-card
            title={article.title ?? article.verbatim ?? 'Untitled article'}
            subtitle={article.journal ?? 'Unknown journal'}
            class="small-card-title bg-surface-light"
            density="compact"
            max-width={600}
          >
            {{
              append: () => <v-chip label text={article.year.toString()} />,
              default: () => <v-card-text>{article.authors.join(', ')}</v-card-text>,
              actions: () =>
                article.doi ? (
                  <v-card-actions>
                    {article.doi ? (
                      <a href={Article.linkDOI(article)}>{Article.linkDOI(article)}</a>
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

export default ArticleChip
