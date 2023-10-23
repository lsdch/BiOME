import useVuelidate from "@vuelidate/core"
import * as spec from "../swagger.json"

import { email, minLength } from "@vuelidate/validators"
import { ValidationArgs } from "@vuelidate/core"

useVuelidate()

export const schemas = spec.definitions

interface Schema {
  properties: object
}

const validators: ValidationArgs = {

}

function parseProperties(schema: Schema) {
  Object.entries(schema.properties).map(([k, v]) => {
    const validators: ValidationArgs = {}
    if ('minLength' in v) {
      validators. = minLength(v.minLength)
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