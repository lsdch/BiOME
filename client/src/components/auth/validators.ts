import zxcvbn from "zxcvbn";

export function passwordValidator(strength: number) {
  return function (password: string, userInputs?: string[]) {
    return zxcvbn(password, userInputs).score >= strength
  }
}