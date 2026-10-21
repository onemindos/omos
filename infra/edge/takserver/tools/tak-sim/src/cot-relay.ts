// Sends CoT XML to FreeTAKServer (TCP) or directly via UDP multicast to ATAK clients
// Also mirrors all CoT to NATS for internal pipeline consumption

import * as net from 'net'
import * as dgram from 'dgram'
import type { NatsConnection } from 'nats'

export interface RelayConfig {
  host: string
  port: number
  protocol: 'tcp' | 'udp'
  nats?: NatsConnection   // optional — mirrors CoT to nats subject
  natsCotSubject?: string
}

export class CotRelay {
  private config: RelayConfig
  private socket: net.Socket | null = null
  private connected = false
  private queue: string[] = []
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null

  constructor(config: RelayConfig) {
    this.config = config
  }

  async connect(): Promise<void> {
    if (this.config.protocol === 'udp') {
      this.connected = true
      return
    }
    return new Promise((resolve) => {
      const sock = new net.Socket()
      sock.connect(this.config.port, this.config.host, () => {
        this.socket = sock
        this.connected = true
        console.log(`[CotRelay] TCP connected → ${this.config.host}:${this.config.port}`)
        this.flushQueue()
        resolve()
      })
      sock.on('error', (err) => {
        console.warn(`[CotRelay] TCP error: ${err.message} — queuing packets, retrying in 5s`)
        this.connected = false
        this.socket = null
        if (!this.reconnectTimer) {
          this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null
            this.connect().catch(() => {})
          }, 5000)
        }
        resolve()  // resolve so startup doesn't block — sends will queue
      })
      sock.on('close', () => {
        this.connected = false
        this.socket = null
        if (!this.reconnectTimer) {
          this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null
            this.connect().catch(() => {})
          }, 5000)
        }
      })
    })
  }

  async send(xml: string): Promise<void> {
    // Mirror to NATS always (for internal pipeline)
    if (this.config.nats && this.config.natsCotSubject) {
      try {
        const sc = this.config.nats
        const subject = this.config.natsCotSubject
        sc.publish(subject, Buffer.from(xml))
      } catch { /* non-fatal */ }
    }

    if (this.config.protocol === 'tcp') {
      if (!this.connected || !this.socket) {
        this.queue.push(xml)
        return
      }
      await new Promise<void>((resolve) => {
        this.socket!.write(xml + '\n', () => resolve())
      })
      return
    }

    // UDP
    await new Promise<void>((resolve) => {
      const client = dgram.createSocket('udp4')
      const buf = Buffer.from(xml)
      client.send(buf, this.config.port, this.config.host, () => {
        client.close()
        resolve()
      })
    })
  }

  private flushQueue(): void {
    const pending = [...this.queue]
    this.queue = []
    for (const xml of pending) {
      this.send(xml).catch(() => {})
    }
  }

  disconnect(): void {
    this.connected = false
    this.socket?.destroy()
    this.socket = null
  }
}
