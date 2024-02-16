import { PeopleService } from "../services/PeopleService"
import { person } from "./fixtures"
import { generateTest } from "./tests"

import e from "../../db/edgeql-js"
import { db } from "./tests"

generateTest("Person", {
  CRUD: {
    list: PeopleService.getPeoplePersons,
    create: PeopleService.createperson,
    update: PeopleService.updatePerson,
    delete: PeopleService.deletePerson,
  },
  getItemIdentifier: ({ id }) => id,
  data: person,
  setup: {
    mockInput: {
      first_name: "Anon",
      last_name: "Ymous"
    },
    async create() {
      return await e.select(
        e.insert(e.people.Person, this.mockInput),
        () => ({
          ...e.people.Person['*'],
          institutions: () => ({
            ...e.people.Institution['*']
          })
        })
      ).run(db)
    },
    async delete(item) {
      return await e.delete(e.people.Person,
        () => ({ filter_single: { id: item.id } })
      ).run(db)
    }
  }
})