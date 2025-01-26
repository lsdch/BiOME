import { ProgramInner } from '@/api'
import { ComponentProps } from 'vue-component-type-helpers'
import { VChip } from 'vuetify/components'

type ProgramChipProps = ComponentProps<typeof VChip> & {
  program: ProgramInner
}

export function ProgramChip({ program, ...props }: ProgramChipProps) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props: menuProps }: { props: any }) => (
          <v-chip text={program.label} {...{ ...menuProps, ...props }} />
        ),
        default: () => (
          <v-card
            max-width={400}
            title={program.label}
            subtitle={program.code}
            class="bg-surface-light"
            density="compact"
          >
            <v-card-text>{program.description}</v-card-text>
          </v-card>
        )
      }}
    </v-menu>
  )
}
