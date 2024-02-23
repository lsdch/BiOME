import { describe, test, expectTypeOf } from "vitest"
import { AuthService } from "../services/AuthService"

import { setupConnection } from "./tests"

setupConnection()

describe("Accounts", () => {
  test("login", async () => {
    await expectTypeOf(AuthService.login({
      identifier: "dev.admin",
      password: "dev.admin",
      remember: true
    })).resolves.toMatchTypeOf({ token: "string" })
  })
  test("login with email", async () => {
    await expectTypeOf(AuthService.login({
      identifier: "dev.admin@mockemail.com",
      password: "dev.admin",
      remember: true
    })).resolves.toMatchTypeOf({ token: "string" })
  })
})