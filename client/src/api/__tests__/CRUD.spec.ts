import { describe, expect, expectTypeOf, test } from 'vitest'

import * as edgedb from "edgedb"
import {
  ApiError,
  CancelablePromise,
  InputValidationError,
  OpenAPI, PeopleService
} from ".."
import { institution, person } from "./fixtures"

OpenAPI.BASE = "http://localhost:5173/api/v1"

const db = edgedb.createClient()

expect.extend({
  toHaveResponseCode(received: ApiError, status: number) {
    if (received.status !== status) return {
      message: () => `expected response status to be ${status}, got ${received.status}`,
      pass: false,
    }
    else return {
      message: () => `response is ${status}`,
      pass: true
    }
  },
  toRespondWithValidationErrors(received: ApiError) {
    if (expectTypeOf(received.body).toMatchTypeOf<InputValidationError[]>()) {
      return {
        message: () => `received validation errors`,
        pass: true
      }
    } else {
      return {
        message: () => `did not receive validation errors`,
        pass: false
      }
    }
  }
})

import { beforeEach } from 'vitest'
import e from "../../db/edgeql-js"

type CRUD<Item, ItemInput, ItemUpdate> = {
  list(): CancelablePromise<Item[]>
  create(input: ItemInput): CancelablePromise<Item>
  update(id: string, input: ItemUpdate): CancelablePromise<Item>
  delete(id: string): CancelablePromise<Item>
}

type TestData<ItemInput, ItemUpdate> = {
  create: ItemInput
  update: ItemUpdate
  invalidUpdate: ItemUpdate
}




type TestEntityConfig<Item, ItemInput, ItemUpdate, MockItem> = {
  CRUD: CRUD<Item, ItemInput, ItemUpdate>,
  getItemIdentifier(item: MockItem | Item): string,
  data: TestData<ItemInput, ItemUpdate>,
  setup: {
    create(): Promise<MockItem>,
    delete(item: MockItem): Promise<any>
  }
}

async function generateTest<
  Item extends MockItem,
  ItemInput extends {},
  ItemUpdate extends {},
  MockItem extends ItemInput & { id: string },
//   // P extends { [k in keyof ObjectTypePointers]: ObjectTypePointers[k] },
//   P extends { id: $BaseObjectλShape["id"] },
//   Z extends ExclusiveTuple,
//   X extends ObjectType<string, $BaseObjectλShape, any, Z>,
//   // {
//   //   [k in keyof ObjectType<string, $SiteλShape, null, $Site['__exclusives__']>]:
//   //   ObjectType<string, $SiteλShape, null, $Site['__exclusives__']>[k]
//   // },
//   // X extends { [k in keyof $BaseObject]: $BaseObject[k] },
//   // T extends { [k in keyof TypeSet<X, edgedb.$.Cardinality.Many>]: TypeSet<X, edgedb.$.Cardinality.Many>[k] },
//   T extends TypeSet<X, edgedb.$.Cardinality.Many>,
//   O extends anonymizeObject<$BaseObject>,
//   Schema extends $expr_PathNode<TypeSet<anonymizeObject<$BaseObject>>>
// // { [k in keyof $expr_PathNode<T, null>]: $expr_PathNode<T, null>[k] }
// // $expr_PathNode<TypeSet<ObjectType<string, X>, edgedb.$.Cardinality.Many>, null>
>(
  entityName: string,
  { CRUD, setup, data, getItemIdentifier }: TestEntityConfig<Item, ItemInput, ItemUpdate, MockItem>
  // schema: $expr_PathNode<ObjectTypeExpression>,
  // schema: Schema,
  // schema: $expr_PathNode<TypeSet<$Institution, edgedb.$.Cardinality.Many>>,
  // setup: InsertShape<O>,
) {
  // const s: (any) => exclusivesToFilterSingle<Schema['__element__']['__exclusives__']> = ({ id }) => id

  // beforeEach(async () => {
  //   const query = e.select(e.insert(schema, setup))
  //   const item: $infer<typeof query> = await query.run(db)

  //   return async () => {
  //     await e.delete(schema, () => ({ filter_single: { id: item.id } })).run(db)
  //   }
  // })

  interface LocalTestContext {
    mock: MockItem
    mockID: string
  }

  return describe<LocalTestContext>(entityName, () => {


    function withMock(...args: Parameters<typeof test<LocalTestContext>>) {
      return test<LocalTestContext>(...args)
    }

    beforeEach<LocalTestContext>(async (context) => {
      context.mock = await setup.create()
      context.mockID = getItemIdentifier(context.mock)
      return async () => {
        await setup.delete(context.mock)
      }
    })

    withMock("lists", async () => {
      await expectTypeOf(CRUD.list()).resolves.toMatchTypeOf<Item[]>()
    })
    withMock("creates", async () => {
      const item = CRUD.create(data.create)
      await expect(item.then((item: Item) => {
        setup.delete(item as MockItem)
        return item
      })).resolves.toMatchObject(data.create)
    })

    withMock('fails to create duplicate', async ({ mock }) => {
      await expect(CRUD.create(mock))
        .rejects.toHaveResponseCode(400)
    })
    withMock("updates", async ({ mockID }) => {
      await expect(CRUD.update(mockID, data.update))
        .resolves.toMatchObject(data.update)
    })
    withMock("fails to update non-existing item", async () => {
      await expect(CRUD.update("null", data.update))
        .rejects.toHaveResponseCode(404)
    })
    withMock("fails to update with invalid property", async ({ mockID }) => {
      await expect(CRUD.update(mockID, data.invalidUpdate))
        .rejects.toHaveResponseCode(400)
    })
    withMock("responds with validation errors on invalid input", async ({ mockID }) => {
      await expect(CRUD.update(mockID, data.invalidUpdate))
        .rejects.toRespondWithValidationErrors()
    })

    withMock("deletes", async ({ mock, mockID }) => {
      await expect(CRUD.delete(mockID))
        .resolves.toMatchObject(mock)
    })
    withMock("fails to delete non-existing item", async () => {
      await expect(CRUD.delete(null))
        .rejects.toHaveResponseCode(404)
    })
  })
}


const entityTests = {
  Institution: {
    CRUD: {
      list: PeopleService.listInstitutions,
      create: PeopleService.createInstitution,
      update: PeopleService.updateInstitution,
      delete: PeopleService.deleteInstitution,
    },
    getItemIdentifier: ({ code }) => code,
    data: institution,
    setup: {
      async create() {
        return await e.select(
          e.insert(e.people.Institution, {
            name: "Vitest institution",
            code: "VITEST"
          }),
          () => ({ ...e.people.Institution['*'] })
        ).run(db)
      },
      async delete(item) {
        return await e.delete(e.people.Institution,
          () => ({ filter_single: { id: item.id } })
        ).run(db)
      }
    }
  },
  Person: {
    CRUD: {
      list: PeopleService.getPeoplePersons,
      create: PeopleService.createperson,
      update: PeopleService.updatePerson,
      delete: PeopleService.deletePerson,
    },
    getItemIdentifier: ({ id }) => id,
    data: person,
    setup: {
      async create() {
        return await e.select(
          e.insert(e.people.Person, {
            first_name: "Anon",
            last_name: "Ymous"
          }),
          () => ({ ...e.people.Person['*'] })
        ).run(db)
      },
      async delete(item) {
        return await e.delete(e.people.Person,
          () => ({ filter_single: { id: item.id } })
        ).run(db)
      }
    }
  }
}

Object.entries(entityTests).forEach(([entity, config]) => {
  generateTest(entity, config)
})
// async () => {
//   const item = await e.select(e.insert(e.people.Institution, {
//     name: "Vitest institution",
//     code: "VITEST"
//   })).run(db)

//   return async () => {
//     await e.delete(e.people.Institution, () => ({ filter_single: { id: item.id } })).run(db)
//   }
// })

// await generateTest('Person', {
//   list: PeopleService.getPeoplePersons,
//   create: PeopleService.createperson,
//   update: PeopleService.updatePerson,
//   delete: PeopleService.deletePerson,
//   getID: ({ id }) => id
// }, person)