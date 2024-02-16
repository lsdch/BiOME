import { PeopleService } from "../services/PeopleService"
import { institution } from "./fixtures"
import { db, generateTest } from "./tests"

import e from "../../db/edgeql-js"

generateTest("Institution", {
  CRUD: {
    list: PeopleService.listInstitutions,
    create: PeopleService.createInstitution,
    update: PeopleService.updateInstitution,
    delete: PeopleService.deleteInstitution,
  },
  getItemIdentifier: ({ code }) => code,
  data: institution,
  setup: {
    mockInput: {
      name: "Vitest institution",
      code: "VITEST"
    },
    async create() {
      return await e.select(
        e.insert(e.people.Institution, this.mockInput),
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