# Biome Frontend

## Production Docker image

From the repository root:

```sh
docker build -f client/docker/Dockerfile -t biome-client:local .
```

The image builds the frontend using Node 24 and `pnpm build` (including
TypeScript checks), then serves the files with Nginx on port 80.
Vue routes can be accessed directly, including after a page reload.
The build uses the generated API client included in the repository.

For Kubernetes, set the container's `BACKEND_URL` environment variable to the
backend Service address and port, for example `http://biome-backend:8080`.
The default is `http://backend:8080`. Use a URL without a path or trailing
slash. The backend Service must exist and its name must be resolvable when
Nginx starts.

Nginx forwards `/api` requests and `/assets` files missing from the frontend to
this Service. Compiled Vite files under `/assets` are served locally.
The Ingress can therefore forward all requests for the domain to the frontend
Service; the API keeps its current URLs on the same domain.
TLS termination can be handled by the Ingress.

To publish the image to the registry used by the cluster:

```sh
docker build -f client/docker/Dockerfile -t <registry>/biome-client:<version> .
docker push <registry>/biome-client:<version>
```

For a local test with the backend already accessible on a Docker network:

```sh
docker run --rm --name biome-client --network biome \
  -p 8081:80 -e BACKEND_URL=http://backend:8080 biome-client:local
```

Adjust the network and backend name, then open `http://localhost:8081`.
The existing Compose configuration is unchanged.
