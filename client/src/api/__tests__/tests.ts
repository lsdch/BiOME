import { describe, expect, expectTypeOf, test } from 'vitest'

import * as edgedb from "edgedb"
import {
  ApiError,
  CancelablePromise,
  InputValidationError,
  Meta,
  OpenAPI
} from ".."

export function setupConnection() {
  OpenAPI.BASE = "http://localhost:8080/api/v1"
  return edgedb.createClient()
}

export const db = setupConnection()

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
  mocks: ItemInput[]
  creates: ItemInput[]
  invalidCreates: ItemInput[]
  update: ItemUpdate
  invalidUpdates: ItemUpdate[]
}

export type TestEntityConfig<Item extends { id: string; meta?: Meta }, ItemInput, ItemUpdate> = {
  CRUD: CRUD<Item, ItemInput, ItemUpdate>,
  getItemIdentifier(item: Item): string,
  data: TestData<ItemInput, ItemUpdate>,
  setup: {
    create(item: ItemInput): Promise<Item>,
    delete(item: Item): Promise<any>
  }
}

export async function generateTest<
  Item extends { id: string; meta?: Meta },
  ItemInput extends {},
  ItemUpdate extends {},
>(
  entityName: string,
  { CRUD, setup, data, getItemIdentifier }: TestEntityConfig<Item, ItemInput, ItemUpdate>
) {


  interface LocalTestContext {
    mockItems: Item[]
    oneMock: Item
    oneID: string
  }

  describe<LocalTestContext>(entityName, () => {


    function withMock(...args: Parameters<typeof test<LocalTestContext>>) {
      return test<LocalTestContext>(...args)
    }

    beforeEach<LocalTestContext>(async (context) => {
      context.mockItems = await Promise.all(data.mocks.map(async (mock) => {
        const item = await setup.create(mock)
        delete item.meta
        return item
      }))
      context.oneMock = context.mockItems[0]
      context.oneID = getItemIdentifier(context.oneMock)

      return async () => {
        await Promise.all(context.mockItems.map(setup.delete))
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
          const item = CRUD.create(data.mocks[0])
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




    withMock("updates", async ({ oneID }) => {
      await expect(CRUD.update(oneID, data.update))
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
          test<LocalTestContext>("responds with BadRequest", async ({ oneID }) => {
            await expect(CRUD.update(oneID, item))
              .rejects.toHaveErrorCode(400)
          })
          test<LocalTestContext>("responds with validation errors on invalid input", async ({ oneID }) => {
            await expect(CRUD.update(oneID, item))
              .rejects.toRespondWithValidationErrors()
          })
        })
    })


    withMock("deletes", async ({ oneMock, oneID }) => {
      await expect(CRUD.delete(oneID))
        .resolves.toMatchObject(oneMock)
    })
    withMock("fails to delete non-existing item", async () => {
      await expect(CRUD.delete("null"))
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