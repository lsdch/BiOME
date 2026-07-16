import { Identification, Taxon } from '@/api'
import { FTaxonStatusIndicator } from '@/features/taxonomy/components/functionals'
import LinkIconGBIF from '@/features/taxonomy/components/LinkIconGBIF'
import { withModifiers } from 'vue'
import { VChip } from 'vuetify/components'

export type IdentificationChipProps = VChip['$props'] & {
  identification: Identification
  short?: boolean
}

export function IdentificationChip({
  identification: { taxon, confer, addendum },
  short,
  ...chipProps
}: IdentificationChipProps) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip
            text={
              short
                ? Taxon.shortName(taxon.name)
                : Identification.taxonString({ taxon, confer, addendum })
            }
            {...{ ...props, ...chipProps }}
            onClick={withModifiers(() => {}, ['stop'])}
          >
            {{
              prepend: () => {
                return confer ? (
                  <v-icon size="small" class="mr-1" color="warning">
                    mdi-tilde
                  </v-icon>
                ) : (
                  <FTaxonStatusIndicator status={taxon.status} />
                )
              }
            }}
          </v-chip>
        ),
        default: () => (
          <v-card
            title={taxon.name}
            subtitle={taxon.authorship}
            class="bg-surface-light small-card-title"
            density="compact"
            to={{ name: 'taxonomy', hash: `#${taxon.name}` }}
          >
            {{
              prepend: () =>
                taxon.gbif_id ? (
                  <LinkIconGBIF
                    GBIF_ID={taxon.gbif_id}
                    variant="tonal"
                    size="x-small"
                    onClick={withModifiers(() => {}, ['stop'])}
                  />
                ) : (
                  <FTaxonStatusIndicator status={taxon.status} />
                ),
              default: () => (
                <div>
                  <em v-if={addendum} class="mx-5">
                    {addendum}
                  </em>
                  <v-card-text>
                    <div class="d-flex justify-space-between">
                      <v-chip text={taxon.status} class="ma-1" size="small" />
                      <v-chip text={taxon.rank} class="ma-1" size="small" />
                    </div>
                  </v-card-text>
                  <v-divider v-if={confer}></v-divider>
                </div>
              ),
              actions: confer
                ? () => (
                    <div>
                      <v-icon size="small" class="mx-2" color="warning">
                        mdi-tilde
                      </v-icon>
                      <span class="text--secondary text-caption">Confer identification</span>
                    </div>
                  )
                : null
            }}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default IdentificationChip
