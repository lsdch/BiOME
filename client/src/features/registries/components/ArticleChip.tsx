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
          <v-chip
            text={Article.toString(article)}
            {...{ ...props, ...chipProps }}
            color={article.original_source ? 'success' : undefined}
          />
        ),
        default: () => (
          <v-card
            title={article.title ?? 'Untitled'}
            subtitle={article.journal ?? 'Unknown journal'}
            class="small-card-title bg-surface-light"
            density="compact"
            max-width={600}
          >
            {{
              append: () => <v-chip label text={article.year.toString()} />,
              default: () => <v-card-text>{article.authors.join(', ')}</v-card-text>,
              actions: () =>
                article.doi || article.original_source ? (
                  <v-card-actions>
                    {article.original_source ? (
                      <v-chip
                        label
                        title="This is the original reference reporting the occurrence"
                        text="Original source"
                      />
                    ) : null}

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
