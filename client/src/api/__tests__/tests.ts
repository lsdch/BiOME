import { describe, expect, expectTypeOf, test } from 'vitest'

import * as edgedb from "edgedb"
import {
  ApiError,
  CancelablePromise,
  InputValidationError,
  OpenAPI
} from ".."

OpenAPI.BASE = "http://localhost:5173/api/v1"

export const db = edgedb.createClient()

expect.extend({
  toHaveErrorCode(received: ApiError, status: number) {
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

export type CRUD<Item, ItemInput, ItemUpdate> = {
  list(): CancelablePromise<Item[]>
  create(input: ItemInput): CancelablePromise<Item>
  update(id: string, input: ItemUpdate): CancelablePromise<Item>
  delete(id: string): CancelablePromise<Item>
}

export type TestData<ItemInput, ItemUpdate> = {
  creates: ItemInput[]
  invalidCreates: ItemInput[]
  update: ItemUpdate
  invalidUpdates: ItemUpdate[]
}




export type TestEntityConfig<Item, ItemInput, ItemUpdate> = {
  CRUD: CRUD<Item, ItemInput, ItemUpdate>,
  getItemIdentifier(item: Item): string,
  data: TestData<ItemInput, ItemUpdate>,
  setup: {
    mockInput: ItemInput,
    create(): Promise<Item>,
    delete(item: Item): Promise<any>
  }
}

export async function generateTest<
  Item extends { id: string },
  ItemInput extends {},
  ItemUpdate extends {},
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
  { CRUD, setup, data, getItemIdentifier }: TestEntityConfig<Item, ItemInput, ItemUpdate>
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
    mock: Item
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

    describe(
      "creates", () => {
        test.each(data.creates)("item %o",
          async (input: ItemInput) => {
            const item = CRUD.create(input)
            await expect(item.then((item: Item) => {
              setup.delete(item)
              return item
            })).resolves
          })
      }
    )

    describe(
      "fails to create", () => {
        withMock('duplicate', async () => {
          const item = CRUD.create(setup.mockInput)
          await expect(item.then((item: Item) => {
            setup.delete(item)
            return item
          }))
            .rejects.toHaveErrorCode(400)
        })
        describe.each(data.invalidCreates)(
          "item %o", (input: ItemInput) => {
            test("responds with BadRequest", async () => {
              const item = CRUD.create(input)
              await expect(item.then((item: Item) => {
                setup.delete(item)
                return item
              }))
                .rejects.toHaveErrorCode(400)
            })
            test("responds with validation errors on invalid input", async () => {
              const item = CRUD.create(input)
              await expect(item.then((item: Item) => {
                setup.delete(item)
                return item
              }))
                .rejects.toRespondWithValidationErrors()
            })
          }
        )
      }
    )




    withMock("updates", async ({ mockID }) => {
      await expect(CRUD.update(mockID, data.update))
        .resolves.toMatchObject(data.update)
    })

    describe("fails to update", () => {
      withMock("non-existing item", async () => {
        await expect(CRUD.update("null", data.update))
          .rejects.toHaveErrorCode(404)
      })

      describe.each(data.invalidUpdates)(
        "with invalid properties %o",
        (item: ItemUpdate) => {
          test<LocalTestContext>("responds with BadRequest", async ({ mockID }) => {
            await expect(CRUD.update(mockID, item))
              .rejects.toHaveErrorCode(400)
          })
          test<LocalTestContext>("responds with validation errors on invalid input", async ({ mockID }) => {
            await expect(CRUD.update(mockID, item))
              .rejects.toRespondWithValidationErrors()
          })
        })
    })


    withMock("deletes", async ({ mock, mockID }) => {
      await expect(CRUD.delete(mockID))
        .resolves.toMatchObject(mock)
    })
    withMock("fails to delete non-existing item", async () => {
      await expect(CRUD.delete(null))
        .rejects.toHaveErrorCode(404)
    })
  })
}



// Object.entries(entityTests).forEach(([entity, config]) => {
//   generateTest(entity, config)
// })
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