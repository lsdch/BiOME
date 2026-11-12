FROM node:24-bookworm-slim AS build

WORKDIR /app/client

RUN npm install --global pnpm@11.20.0

ENV CYPRESS_INSTALL_BINARY=0

COPY client/package.json client/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --ignore-scripts

COPY client/ ./
COPY biome_logo.svg /app/biome_logo.svg
RUN pnpm build

FROM nginx:stable-alpine AS runtime

ENV BACKEND_URL=http://backend:8080

COPY client/docker/default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/client/dist /usr/share/nginx/html

EXPOSE 80
