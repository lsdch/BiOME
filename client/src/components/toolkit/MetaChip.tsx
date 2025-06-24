import { Meta } from '@/api'
import { DateTime } from 'luxon'
import { VChip } from 'vuetify/components'

export type MetaChipProps = VChip['$props'] & {
  meta: Meta
  iconColor?: string
}

export function MetaChip({ meta, iconColor = '#777', ...chipProps }: MetaChipProps) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip
            label
            color={iconColor}
            variant="tonal"
            text={Meta.toString(meta)}
            {...{ ...props, ...chipProps }}
          >
            {{
              prepend: () => (
                <v-icon color={iconColor} class="mr-3" size="small" icon={Meta.icon(meta)} />
              )
            }}
          </v-chip>
        ),
        default: () => (
          <v-card min-width={350} class="bg-surface-light" density="compact">
            <v-list class="bg-surface-light" density="compact">
              <v-list-item
                title={DateTime.fromJSDate(meta.created).toLocaleString(DateTime.DATE_FULL, {
                  locale: 'en-gb'
                })}
                subtitle={`Created at ${DateTime.fromJSDate(meta.created).toLocaleString(
                  DateTime.TIME_24_SIMPLE,
                  { locale: 'en-gb' }
                )}`}
                prepend-icon="mdi-content-save"
                slim
              >
                {{
                  append: () => <v-chip text={meta.created_by?.name ?? 'No author'} size="small" />
                }}
              </v-list-item>
              <v-divider />
              <v-list-item
                title={
                  meta.modified
                    ? DateTime.fromJSDate(meta.modified).toLocaleString(DateTime.DATE_FULL, {
                        locale: 'en-gb'
                      })
                    : undefined
                }
                subtitle={
                  meta.modified
                    ? `Last updated at ${DateTime.fromJSDate(meta.modified).toLocaleString(
                        DateTime.TIME_24_SIMPLE,
                        { locale: 'en-gb' }
                      )}`
                    : 'Never updated'
                }
                prepend-icon="mdi-update"
                slim
              />
            </v-list>
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default MetaChip
