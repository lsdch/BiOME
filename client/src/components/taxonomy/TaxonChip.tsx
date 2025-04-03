import { Taxon } from '@/api'
import { FTaxonStatusIndicator } from '@/components/taxonomy/functionals'
import LinkIconGBIF from '@/components/taxonomy/LinkIconGBIF'
import { withModifiers } from 'vue'
import { VChip } from 'vuetify/components'

export type TaxonChipProps = VChip['$props'] & { taxon: Taxon; short?: boolean }

export function TaxonChip({ taxon, short, ...chipProps }: TaxonChipProps) {
  const shortName = (name: string) => {
    const splitName = name.split(' ')
    if (splitName.length === 1) return name
    else return `${splitName[0][0]}. ${splitName.slice(1).join(' ')}`
  }

  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip
            text={short ? shortName(taxon.name) : taxon.name}
            {...{ ...props, ...chipProps }}
            onClick={withModifiers(() => {}, ['stop'])}
          />
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
                taxon.GBIF_ID ? (
                  <LinkIconGBIF
                    GBIF_ID={taxon.GBIF_ID}
                    variant="tonal"
                    size="x-small"
                    onClick={withModifiers(() => {}, ['stop'])}
                  />
                ) : (
                  <FTaxonStatusIndicator status={taxon.status} />
                ),
              default: () => (
                <v-card-text>
                  <div class="d-flex justify-space-between">
                    <v-chip text={taxon.status} class="ma-1" size="small" />
                    <v-chip text={taxon.rank} class="ma-1" size="small" />
                  </div>
                </v-card-text>
              )
            }}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default TaxonChip
