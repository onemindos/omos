FROM oven/bun:1 AS builder

RUN apt-get update && apt-get install -y git ca-certificates && rm -rf /var/lib/apt/lists/*
RUN git clone --depth=1 https://github.com/synadia-ai/synadia-agents.git /synadia

WORKDIR /synadia/examples/agent-web-ui
RUN bun install
RUN bun run build

FROM oven/bun:1

WORKDIR /app

COPY --from=builder /synadia/examples/agent-web-ui/dist ./dist
COPY --from=builder /synadia/examples/agent-web-ui/server/bridge.ts ./server/bridge.ts
COPY --from=builder /synadia/examples/agent-web-ui/server/wire.ts ./server/wire.ts
COPY --from=builder /synadia/examples/agent-web-ui/package.json ./package.json
COPY server-index.ts ./server/index.ts
RUN bun install --production

EXPOSE 3300
CMD ["bun", "run", "./server/index.ts"]
