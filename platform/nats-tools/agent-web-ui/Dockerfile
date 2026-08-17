# agent-web-ui has no upstream Dockerfile — lives inside synadia-ai/synadia-agents
# monorepo at examples/agent-web-ui/. The SDK dependency (@synadia-ai/agents) is a
# file:../../client-sdk/typescript path ref in package.json — we resolve it by
# building the SDK first from the monorepo root before installing example deps.
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy full monorepo so the local file: SDK ref resolves
COPY . .

# Build the SDK (client-sdk/typescript) so the file: dep resolves
RUN cd client-sdk/typescript && bun install && bun run build

# Install agent-web-ui deps (resolves file:../../client-sdk/typescript)
RUN cd examples/agent-web-ui && bun install

# Build frontend assets
RUN cd examples/agent-web-ui && bun run build

FROM oven/bun:1-alpine AS runner

WORKDIR /app

# Copy monorepo slices needed at runtime
COPY --from=builder /app/client-sdk/typescript/dist ./client-sdk/typescript/dist
COPY --from=builder /app/client-sdk/typescript/package.json ./client-sdk/typescript/package.json
COPY --from=builder /app/examples/agent-web-ui ./examples/agent-web-ui

WORKDIR /app/examples/agent-web-ui

ENV NODE_ENV=production
EXPOSE 3300

CMD ["bun", "run", "server/index.ts"]
