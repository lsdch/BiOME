import { PeopleService } from "../services/PeopleService"
import { institution } from "./fixtures"
import { TestData, db, generateTest } from "./tests"

import e from "../../db/edgeql-js"
import { InstitutionInput } from "../models/InstitutionInput"
import { InstitutionUpdate } from "../models/InstitutionUpdate"
import { Institution } from "../models/Institution"

generateTest("Institution", {
  CRUD: {
    list: PeopleService.listInstitutions,
    create: PeopleService.createInstitution,
    update: PeopleService.updateInstitution,
    delete: PeopleService.deleteInstitution,
  },
  getItemIdentifier: ({ code }) => code,
  data: <TestData<InstitutionInput, InstitutionUpdate>>institution,
  setup: {
    async create(mockInput: InstitutionInput): Promise<Institution> {
      return await e.select(
        e.insert(e.people.Institution, mockInput),
        () => ({ ...e.people.Institution['*'] })
      ).run(db)
    },
    async delete(item) {
      return await e.delete(e.people.Institution,
        () => ({ filter_single: { id: item.id } })
      ).run(db)
    }
  }
})