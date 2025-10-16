export function CongruenceChip(
  { is_congruent }: { is_congruent: boolean },
  context: { attrs?: object }
) {
  return (
    <v-tooltip
      text={
        is_congruent
          ? 'Bio material identification matches its sequences identification'
          : 'Bio material identification contradicted by its sequences identification'
      }
      open-on-click
      location="end"
      origin="center"
    >
      {{
        activator({ props }: { props: object }) {
          return (
            <v-chip
              {...{
                ...props,
                ...(is_congruent
                  ? {
                      color: 'success',
                      text: 'Congruent'
                    }
                  : {
                      color: 'warning',
                      text: 'Incongruent'
                    })
              }}
              size="small"
            />
          )
        }
      }}
    </v-tooltip>
  )
}

export default CongruenceChip
