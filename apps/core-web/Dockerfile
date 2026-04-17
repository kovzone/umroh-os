# Multi-stage Dockerfile for apps/core-web.
#
# Targets:
#   - deps:  install npm dependencies (cached layer; reused by dev + build)
#   - dev:   development image — Vite dev server with hot reload. Used by
#            docker-compose.dev.yml (target: dev) with the source bind-mounted.
#   - build: produces the SvelteKit/adapter-node output in /app/build/
#   - prod:  minimal Node runtime that runs `node build/index.js`. Not yet
#            wired into compose; reserved for the eventual production deploy.
#
# Build context must be the repo root so the Dockerfile can `COPY apps/core-web ...`.
# e.g.: docker build -f apps/core-web/Dockerfile --target dev .

# ─── Stage 1: deps ────────────────────────────────────────────────────────
FROM node:20-alpine AS deps
WORKDIR /app
COPY apps/core-web/package.json apps/core-web/package-lock.json* ./
RUN npm install

# ─── Stage 2: dev (used by docker-compose.dev.yml) ───────────────────────
FROM node:20-alpine AS dev
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
# Source is bind-mounted by compose; nothing else copied here so changes
# on the host are reflected immediately. node_modules is shadowed by an
# anonymous volume in compose to keep the container's install in place.
EXPOSE 3000
CMD ["npm", "run", "dev"]

# ─── Stage 3: build ──────────────────────────────────────────────────────
FROM node:20-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY apps/core-web ./
RUN npm run build

# ─── Stage 4: prod runtime ───────────────────────────────────────────────
FROM node:20-alpine AS prod
WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app/build ./build
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./
EXPOSE 3000
CMD ["node", "build/index.js"]
