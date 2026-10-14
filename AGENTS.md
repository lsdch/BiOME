# BiOME — Agent Instructions

BiOME is a scientific biodiversity data platform for importing, resolving, storing, querying, and visualizing ecological occurrence and sampling data.

## Repository structure

The project contains:

* A Go backend exposing a REST API.
* A Vue 3 / TypeScript frontend.
* PostgreSQL as the primary database.
* SQL queries and generated Go database code managed by `sqlc` and `go-jet`.
* OpenAPI-generated TypeScript API types/client code.

Before modifying code, inspect the relevant package/module and follow existing patterns. Avoid broad repository exploration unless necessary for the task.

## Backend

Main technologies:

* Go
* Gin
* PostgreSQL
* PostGIS
* H3
* sqlc

### Go conventions

* Follow idiomatic Go and existing project conventions.
* Keep changes focused on the requested task.
* Prefer simple explicit code over unnecessary abstractions.
* Propagate `context.Context` through database, storage, and external API operations.
* Handle errors explicitly.
* Do not silently ignore errors.
* Preserve existing domain types where possible instead of introducing duplicate representations.

Run formatting on modified Go files.

When practical, run relevant tests after changes.

## Database

PostgreSQL is the source of truth for persistent application data.

Relevant extensions include:

* PostGIS
* H3
* pgcrypto
* citext
* pg_trgm

### SQL rules

* Prefer SQL operations over reproducing relational logic in Go.
* Keep queries compatible with `sqlc`.
* Use explicit joins and predicates.
* Preserve database constraints and referential integrity.
* Consider `NULL` semantics explicitly.
* Avoid unnecessary queries and N+1 patterns.
* Use transactions when multiple writes must be atomic.
* Do not modify generated sqlc Go files manually.

When changing SQL consumed by sqlc, regenerate the generated code using the project's existing generation workflow.

### Schema changes

When modifying the database schema:

* Consider foreign-key behavior (`CASCADE`, `RESTRICT`, etc.).
* Add indexes when justified by actual query patterns.
* Avoid redundant indexes already covered by a primary key or another suitable index.
* Preserve existing data semantics.
* Consider migration compatibility rather than assuming an empty database.

## Biodiversity data model

The application handles scientific biodiversity data, including:

* taxa and taxonomic resolution
* occurrences
* samplings
* habitats
* datasets and import batches
* publications
* geographic information

Imported values may intentionally preserve verbatim source data separately from normalized/resolved values.

Do not discard or normalize verbatim scientific data unless explicitly required.

Taxonomic names may be unresolved, ambiguous, qualified, or matched against external sources such as GBIF. Do not assume every imported taxon has an exact canonical match.

Dates may have variable precision (year, month, day) or be unknown. Do not assume all dates represent exact days.

## Spatial data

Spatial functionality uses PostGIS and H3.

* Geographic coordinates use WGS84 / SRID 4326 unless existing code explicitly specifies otherwise.
* Do not confuse H3 indexes with arbitrary integer identifiers.
* Preserve the distinction between precise sampling locations and H3-derived aggregations.
* Avoid expensive spatial operations when an existing indexed PostGIS/H3 operation can perform the work.

## Import pipeline

Imports may involve staging tables and resolution steps before data is materialized into canonical tables.

Preserve the distinction between:

1. raw/verbatim imported data
2. staging data
3. resolution/candidate data
4. canonical persisted data

Do not bypass resolution or staging stages merely to simplify an implementation.

Uploaded source files may also be retained by the application's file-storage layer. Use the existing storage abstraction rather than directly accessing arbitrary filesystem paths.

## Frontend

Main technologies:

* Vue 3
* TypeScript
* Vuetify
* TanStack Query
* Vite
* generated OpenAPI client/types
* MapLibre GL JS
* deck.gl

### Vue conventions

* Use Vue 3 Composition API.
* Follow existing component and composable patterns.
* Keep state local unless it genuinely belongs in shared state.
* Prefer computed state over duplicated synchronized state.
* Preserve TypeScript type safety; avoid `any` unless there is a concrete reason.

### API access

Use the generated OpenAPI client and generated types whenever possible.

Do not manually edit generated OpenAPI client/type files.

Do not duplicate generated API types with handwritten equivalents unless a separate UI/domain representation is genuinely necessary.

When API parameters are optional, preserve the distinction between:

* absent / `undefined`
* `null`
* zero/false/empty values

Do not assume these are interchangeable.

### TanStack Query

Query keys must contain every value that affects the request.

When changing filters or request parameters, verify both:

* the query key changes appropriately
* the serialized HTTP request contains the expected parameters

Avoid unnecessary manual refetching when query-key reactivity can express the dependency.

## Maps

Map rendering uses MapLibre and deck.gl.

Preserve the distinction between:

* exact occurrence/sampling positions
* H3 aggregated representations

Avoid expensive client-side geographic computations when they can reasonably be performed by the backend/database.

## Generated code

Generated files must not be manually edited unless the task explicitly requires changing the generator output itself.

Typical generated code includes:

* sqlc Go code
* OpenAPI TypeScript client/types

Modify the source schema/query/generator input and regenerate instead.

## Working principles

For each task:

1. Identify the smallest relevant part of the codebase.
2. Read the existing implementation before changing it.
3. Follow existing architecture and naming.
4. Make the smallest coherent change that solves the problem.
5. Avoid unrelated refactoring.
6. Update tests when behavior changes.
7. Run the narrowest relevant validation first.
8. Expand validation only when necessary.

Do not introduce a new dependency when the standard library or an existing project dependency already provides an adequate solution.

Do not change public APIs, database schemas, generated interfaces, or persisted data formats unless required by the task.

If a requirement is ambiguous and the choice would significantly affect architecture or persisted data, ask before making that choice.

Never instanciate a backend service outside of the `App` initialisation.