import { describe, it, expect, expectTypeOf } from 'vitest'

import {
  OpenAPI, PeopleService, ApiError, InputValidationError, CancelablePromise
} from ".."
import institution from "./fixtures/people/institution.json"

OpenAPI.BASE = "http://localhost:5173/api/v1"

expect.extend({
  toHaveResponseCode(received: ApiError, status: number) {
    if (received.status !== status) return {
      message: () => `expected response status to be ${status}`,
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

type CRUD<Item, ItemInput, ItemUpdate> = {
  list(): CancelablePromise<Item[]>
  create(input: ItemInput): CancelablePromise<Item>
  update(id: string, input: ItemUpdate): CancelablePromise<Item>
  delete(id: string): CancelablePromise<Item>
  getID(item: ItemInput): string
}

type TestData<ItemInput, ItemUpdate> = {
  create: ItemInput
  update: ItemUpdate
  invalidUpdate: ItemUpdate
}

function generateTest<Item extends ItemInput, ItemInput extends {}, ItemUpdate extends {}>(
  entityName: string,
  crud: CRUD<Item, ItemInput, ItemUpdate>,
  data: TestData<ItemInput, ItemUpdate>
) {
  return describe(entityName, () => {
    it("lists", async () => {
      await expectTypeOf(crud.list()).resolves.toMatchTypeOf<Item[]>()
    })
    it("creates", async () => {
      await expect(crud.create(data.create))
        .resolves.toMatchObject(data.create)
    })
    it('fails to create duplicate', async () => {
      await expect(crud.create(data.create))
        .rejects.toHaveResponseCode(400)
    })
    it("updates", async () => {
      await expect(crud.update(crud.getID(data.create), data.update))
        .resolves.toMatchObject(data.update)
    })
    it("fails to update non-existing item", async () => {
      await expect(crud.update("NOTFOUND", data.update))
        .rejects.toHaveResponseCode(404)
    })
    it("fails to update with invalid property", async () => {
      await expect(crud.update(crud.getID(data.create), data.invalidUpdate))
        .rejects.toHaveResponseCode(400)
    })
    it("responds with validation errors on invalid input", async () => {
      await expect(crud.update(crud.getID(data.create), data.invalidUpdate))
        .rejects.toRespondWithValidationErrors()
    })

    it("deletes", async () => {
      await expect(crud.delete(crud.getID(data.create)))
        .resolves.toMatchObject({ code: crud.getID(data.create) })
    })
    it("fails to delete non-existing item", async () => {
      await expect(crud.delete(crud.getID(data.create)))
        .rejects.toHaveResponseCode(404)
    })
  })
}

generateTest('Institution', {
  list: PeopleService.listInstitutions,
  create: PeopleService.createInstitution,
  update: PeopleService.updateInstitution,
  delete: PeopleService.deleteInstitution,
  getID: ({ code }) => code
}, institution)
