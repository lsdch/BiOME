# Biome Backend

## PostgreSQL migrations

From the repository root, start the development database and apply its migrations:

```sh
docker compose up -d db
task server:db:migrate
task server:db:migration:status
```

The database must already exist and have the extension packages available (including
PostGIS and H3). The tasks default to the local Compose database. Set `DB_URL` in
the environment to target another database. The migration role needs permission
to create the application objects and the `migrations` schema, and to enable any
extensions that have not already been installed in the target database.

The first 33 migrations reproduce the original SQL bootstrap. Start with a fresh
database for a new installation. For a database initialized by the old bootstrap,
use a separate empty database if the data is disposable. If data must be retained,
back up the database and verify its schema against a freshly migrated database
before registering the initial versions as applied. Test that adoption on a
restored copy first.

To change the schema:

```sh
task server:db:migration:create NAME=add_example_column
# Edit the generated SQL migration under server/db/migrations.
task server:db:migration:validate
task server:db:migrate
task server:db:generate
```

Migration filenames use four-digit sequential versions, such as
`0034_add_example_column.sql`. Commit new migrations and any resulting generated
code together. Keep applied migrations immutable. sqlc reads the migration files in filename order, and Go-Jet
reads the migrated `public` schema. Goose tracks applied versions in
`migrations.goose_db_version`. The initial migrations contain only an upward path.

For a deployed instance, run migrations as a single deployment step before starting
the new backend version.

## Docker image 

From the repository root:

```sh
docker build -f server/docker/Dockerfile -t biome-server:local ./server
```

The local `server/huma` module is required by `go.mod` and must be present in
the build context, including in CI. CGO is retained for H3; both the build and
runtime stages use Debian Bookworm. The image includes CA certificates,
time zone data and assets, but no deployment configuration files.

The same image can run a migration Job before the backend rollout. Override
the entrypoint with `/app/goose` and pass these arguments:

```sh
/app/goose -env=none -dir /app/migrations -table migrations.goose_db_version up
```

Set `GOOSE_DRIVER=postgres` and inject `GOOSE_DBSTRING` from the migration
role's Secret. Provision the `migrations` schema for that role before running
the Job. Use the same image digest for the Job and the backend.

## Kubernetes configuration

Mount a `config.yml` file from a ConfigMap read-only into `/etc/biome`.
This file is required. Start from `server/config.yml` and complete the
installation-specific settings. The image supplies these values:

| Variable | Value |
| --- | --- |
| `BIOME_ENV` | `prod` |
| `GIN_MODE` | `release` |
| `BIOME_API_API_HOST` | `0.0.0.0` |
| `BIOME_API_API_PORT` | `8080` |
| `BIOME_RAW_FILE_STORAGE_ROOT` | `/data/raw_files` |
| `BIOME_BOOTSTRAP_COUNTRIES_COUNTRIES_JSON_CACHE_PATH` | `/data/cache/countries.json` |

`BIOME_*` environment variables take precedence over configuration files.
Their names follow the nested YAML keys, separated by `_`; repetitions such
as `API_API_PORT` are intentional. Keys may be absent from the YAML file.

Inject secrets through `env[].valueFrom.secretKeyRef`, particularly:

```yaml
env:
  - name: BIOME_DB_DSN
    valueFrom:
      secretKeyRef:
        name: biome-backend
        key: database-dsn
  - name: BIOME_AUTH_TOKENS_JWT_SECRET_KEY
    valueFrom:
      secretKeyRef:
        name: biome-backend
        key: jwt-secret-key
  - name: BIOME_AUTH_TOKENS_REFRESH_TOKEN_PEPPER
    valueFrom:
      secretKeyRef:
        name: biome-backend
        key: refresh-token-pepper
```

Also provide the public URL (`BIOME_APP_PUBLIC_BASE_URL`), administrator
contact (`BIOME_INSTANCE_ADMIN_EMAIL`), API metadata (`BIOME_API_VERSION`,
`BIOME_API_API_CONTACT_NAME`, `BIOME_API_API_CONTACT_EMAIL`), API base path
(`BIOME_API_API_BASE_PATH`, usually `/api/v1`) and `BIOME_GBIF_GBIF_USER_AGENT`,
either in the configuration file or in the Deployment environment variables.
`JWT_SECRET_KEY` requires at least 32 characters, and `REFRESH_TOKEN_PEPPER`
requires at least 8.

The optional files `/etc/biome/config.prod.yml` and then
`/etc/biome/config.local.yml` override the base configuration, before
environment variables are applied. Configuration changes require a Pod
restart. SMTP settings already persisted in the database take precedence
over configuration initialization values.

The process runs with UID/GID `10001:10001`. Provide writable volumes with
matching permissions (`fsGroup: 10001` if supported by the storage driver):

| Path | Usage |
| --- | --- |
| `/data/raw_files` | Imported source files; shared persistent storage if using multiple replicas |
| `/data/cache` | Country cache; `emptyDir` can be used |
| `/data/openapi` | Specifications generated at startup; `emptyDir` can be used |
| `/app/assets` | Images and editable application icon; persistent storage if needed |

Mounting an empty volume at `/app/assets` hides the files supplied by the
image: initialize its contents before use. `/app/data` points to `/data` for
legacy relative paths. The Kubernetes Service must target port 8080 unless
`BIOME_API_API_PORT` is overridden.

The default command is
`/app/server -config-dir /etc/biome -openapi /data/openapi`. The Deployment's
`args` can override these options. Local image checks do not replace testing
with PostgreSQL and the cluster's actual volumes.
