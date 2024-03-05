import { zxcvbn } from '@zxcvbn-ts/core'

export function passwordValidator(options: { strength: number }) {
  return function (password: string, userInputs?: string[]) {
    return zxcvbn(password, userInputs).score >= options.strength
  }
}