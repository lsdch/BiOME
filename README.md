## Dev environment

### Setup

1. Install golang >= 1.21
2. Setup `GOBIN` environment variable in your shell profile
3. Add it to your `PATH` environment variable
4. Install [air](github.com/cosmtrek/air@latest) to handle hot-reloading of the server, and [swag](github.com/swaggo/swag/cmd/swag@latest) to generate OpenAPI docs :

   ```bash
   go install github.com/cosmtrek/air@latest
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
5. Database setup
  - install [EdgeDB CLI](https://www.edgedb.com/docs/intro/cli)
  - run `edgedb project init` in the `/server` directory

6. Client setup
- Install [pnpm](https://pnpm.io/installation)
- run `pnpm i` in the `/client` directory

### Running the local server

- run `air` in the `/server` directory
- run `pnpm dev` in the `/client` directory
