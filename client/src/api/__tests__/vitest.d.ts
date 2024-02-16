import 'vitest'

interface CustomMatchers<R = unknown> {
  toHaveErrorCode(status: number): R
  toRespondWithValidationErrors(): R
}

declare module 'vitest' {
  interface Assertion<T = any> extends CustomMatchers<T> { }
  interface AsymmetricMatchersContaining extends CustomMatchers { }
}