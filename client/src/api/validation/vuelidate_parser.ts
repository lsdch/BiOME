import useVuelidate from "@vuelidate/core"
import * as spec from "../swagger.json"

import { email, maxLength, minLength, minValue, maxValue, sameAs } from "@vuelidate/validators"
import { ValidationArgs } from "@vuelidate/core"
import { ValidationRule } from "@vuelidate/core"
import { ValidationRuleCollection } from "@vuelidate/core"

useVuelidate()

export const schemas = spec.definitions

interface Schema {
  properties: object
}

type ValidationMappings<T = any> = Record<string, (value: T) => ValidationRule<T>>

const validationMappings: ValidationMappings = {
  minLength: minLength,
  maxLength: maxLength,
  minimum: minValue,
  maximum: maxValue,
}

function parseProperties(schema: Schema) {
  Object.entries(schema.properties).map(([k, v]) => {
    let validators: ValidationArgs = {}
    if ('minLength' in v) {
      validators = { minLength: minLength(v.minLength), ...validators }
    }
    if ('maxLength' in v) {
      validators = { maxLength: maxLength(v.maxLength), ...validators }
    }
    if ('minimum' in v) {
      validators =
    }
  })
}

function parseFormat(format: string): object {
  switch (format) {
    case "email":
      return { email }
      break;

    default:
      return {}
      break;
  }
}