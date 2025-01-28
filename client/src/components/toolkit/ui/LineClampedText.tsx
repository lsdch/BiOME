import { SetupContext } from 'vue'

export function LineClampedText(
  { text, lines }: { text: string; lines: number },
  context: SetupContext
) {
  return (
    <div
      class="line-clamped-text"
      style={{
        overflow: 'hidden',
        display: '-webkit-box',
        lineClamp: lines,
        '-webkit-box-orient': 'vertical'
      }}
      {...context.attrs}
    >
      {text}
    </div>
  )
}
