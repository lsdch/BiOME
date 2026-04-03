# Database queries 
We use Gel as our database manager, so when suggesting database queries, use EdgeQL syntax. When the query is embedded in Go code, always use backticks string delimiters and prefix with #edgeql.
Example:
`#edgeql
{ query goes here }
`

Database queries must use the database schema defined in the files db/schema/*.esdl

# Front end

## Typescript code style

- Use camelCase for variable and function names, and PascalCase for type and interface names.
- Avoid using the `any` type. Instead, define specific types or interfaces for your data structures.
- Use `const` for variables that are not reassigned, and `let` for variables that are reassigned. Avoid using `var`.
- Avoid using for loops. Instead, use array methods like `map`, `filter`, `foreach` and `reduce` for iterating over arrays.

## API client

The API client is automatically generated from the OpenAPI specification defined in client/openapi.json. The OpenAPI specification is generated from the backend code using the Go framework Huma. When suggesting frontend code, do not attempt to change the API client directly in client/src/api, but instead suggest changes to the backend code in server/src, which will then be reflected in the OpenAPI specification and the generated API client.