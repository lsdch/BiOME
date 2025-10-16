import { SeqReference } from '@/api'

export function SeqRefChip({ seqRef }: { seqRef: SeqReference }, context: { attrs?: object }) {
  const { db, accession, is_origin } = seqRef
  return (
    <v-chip
      text={`${db.code}: ${accession}`}
      href={SeqReference.link(seqRef)}
      target="_blank"
      color={is_origin ? 'success' : undefined}
      prepend-icon={is_origin ? 'mdi-arrow-down' : undefined}
      {...context.attrs}
    />
  )
}

export default SeqRefChip
