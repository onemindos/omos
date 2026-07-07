import { join, extname } from "node:path";
import { existsSync, statSync } from "node:fs";
import { Agents, SDK_PROTOCOL_VERSION } from "@synadia-ai/agents";
import { connect as natsConnect } from "@nats-io/transport-node";
import { Bridge, formatSdkProtocolVersion, type BridgeWsData } from "./bridge.ts";

const PORT = parseInt(process.env["PORT"] ?? "3300", 10);
const NATS_URL = process.env["NATS_URL"] ?? "nats://localhost:4222";
const NATS_USER = process.env["NATS_USER"] ?? "";
const NATS_PASSWORD = process.env["NATS_PASSWORD"] ?? "";

const nc = await natsConnect({
  servers: NATS_URL,
  user: NATS_USER,
  pass: NATS_PASSWORD,
  tls: { rejectUnauthorized: false },
  name: "synadia-bridge",
});

const agents = new Agents({ nc });
console.log(`[synadia-bridge] NATS connected: ${NATS_URL}`);

const distDir = join(import.meta.dir, "..", "dist");
const sdkVersionString = formatSdkProtocolVersion(SDK_PROTOCOL_VERSION);

const server = Bun.serve<BridgeWsData>({
  port: PORT,
  async fetch(req, srv) {
    const url = new URL(req.url);
    if (url.pathname === "/ws") {
      const bridge = new Bridge(agents, nc, sdkVersionString);
      const upgraded = srv.upgrade(req, { data: { bridge } });
      if (upgraded) return undefined;
      return new Response("expected WebSocket upgrade on /ws", { status: 400 });
    }
    if (existsSync(distDir)) {
      const safePath = url.pathname === "/" ? "/index.html" : url.pathname;
      const filePath = join(distDir, decodeURIComponent(safePath));
      if (!filePath.startsWith(distDir)) return new Response("forbidden", { status: 403 });
      if (existsSync(filePath) && statSync(filePath).isFile()) return new Response(Bun.file(filePath));
      if (!extname(url.pathname)) return new Response(Bun.file(join(distDir, "index.html")));
      return new Response("not found", { status: 404 });
    }
    return new Response("dist not found\n", { status: 404 });
  },
  websocket: {
    open(ws) { ws.data.bridge.open(ws); },
    message(ws, msg) {
      const text = typeof msg === "string" ? msg : new TextDecoder().decode(msg);
      ws.data.bridge.onMessage(text);
    },
    close(ws) { ws.data.bridge.close(); },
  },
});

console.log(`[synadia-bridge] listening on http://localhost:${server.port}`);

async function shutdown(sig: string): Promise<void> {
  console.log(`[synadia-bridge] ${sig} — shutting down`);
  try { server.stop(); } catch { /* noop */ }
  try { await agents.close(); } catch { /* noop */ }
  try { await nc.close(); } catch { /* noop */ }
  process.exit(0);
}

process.on("SIGINT", () => void shutdown("SIGINT"));
process.on("SIGTERM", () => void shutdown("SIGTERM"));
