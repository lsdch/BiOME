import { SeqReference } from '@/api'

export function SeqRefChip({ seqRef }: { seqRef: SeqReference }, context: { attrs?: object }) {
  const { db, accession, is_origin } = seqRef
  return (
    <v-tooltip>
      {{
        default: () =>
          is_origin ? (
            <span>
              The <span class="text-success font-weight-bold">original source</span> of the
              sequence{' '}
            </span>
          ) : (
            'Sequence was deposited here'
          ),
        activator: ({ props }: any) => (
          <v-chip
            text={`${db.code}:${accession}`}
            href={SeqReference.link(seqRef)}
            target="_blank"
            color={is_origin ? 'success' : undefined}
            prepend-icon={is_origin ? 'mdi-arrow-down' : undefined}
            title={is_origin ? 'Sequence origin' : 'External sequence reference'}
            {...{ ...props, ...context.attrs }}
          />
        )
      }}
    </v-tooltip>
  )
}

export default SeqRefChip
