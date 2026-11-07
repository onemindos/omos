import type { LoggerService } from '@backstage/backend-plugin-api';
import { connect, type NatsConnection, type Subscription } from 'nats';
import type { SourceStatus } from './contract';
import { RateMeter } from './rate';

/** Subject filters the bridge subscribes to. Core subscriptions only; nothing is ever published. */
export const SUBSCRIPTIONS = [
  'tak.cot.>',
  'intel.>',
  'evt.>',
  'det.>',
  'ent.>',
  'tel.>',
  'cmd.>',
  'agents.status.>',
] as const;

export interface NatsSourceOptions {
  url: string;
  user?: string;
  password?: string;
  caFile?: string;
  logger: LoggerService;
  onMessage(subject: string, data: Uint8Array): void;
  now?: () => number;
}

export function ignoredSubject(subject: string): boolean {
  return subject.startsWith('_INBOX.') || subject.startsWith('$');
}

/**
 * Read-only NATS subscriber. Without credentials the state is `disabled` and nothing connects.
 * The client reconnects forever; the initial connection is retried here with backoff.
 */
export class NatsSource {
  private nc?: NatsConnection;
  private stopped = false;
  private state: SourceStatus['state'];
  private detail?: string;
  private lastAt?: number;
  private readonly meter = new RateMeter();
  private readonly now: () => number;

  constructor(private readonly opts: NatsSourceOptions) {
    this.now = opts.now ?? Date.now;
    this.state = opts.user && opts.password ? 'connecting' : 'disabled';
    if (this.state === 'disabled')
      this.detail = 'NATS_USER / NATS_PASSWORD not set';
  }

  start(): void {
    if (this.state === 'disabled') return;
    void this.run();
  }

  async stop(): Promise<void> {
    this.stopped = true;
    try {
      await this.nc?.drain();
    } catch {
      // closing anyway
    }
  }

  status(): SourceStatus {
    return {
      state: this.state,
      lastMessageAt: this.lastAt
        ? new Date(this.lastAt).toISOString()
        : undefined,
      rate: this.meter.rate(this.now()),
      detail: this.detail,
    };
  }

  private async run(): Promise<void> {
    let backoff = 1000;
    while (!this.stopped) {
      try {
        const nc = await connect({
          servers: this.opts.url,
          user: this.opts.user,
          pass: this.opts.password,
          name: 'onemind-portal-livebridge',
          maxReconnectAttempts: -1,
          reconnectTimeWait: 2000,
          timeout: 10_000,
          noEcho: true,
          ...(this.opts.caFile || this.opts.url.startsWith('tls://')
            ? { tls: this.opts.caFile ? { caFile: this.opts.caFile } : {} }
            : {}),
        });
        this.nc = nc;
        backoff = 1000;
        this.state = 'connected';
        this.detail = undefined;
        this.opts.logger.info(`livebridge NATS connected to ${nc.getServer()}`);
        void this.watchStatus(nc);
        this.subscribe(nc);
        await nc.closed();
        if (!this.stopped) {
          this.state = 'down';
          this.detail = 'connection closed';
        }
      } catch (e) {
        this.state = 'down';
        this.detail = String((e as Error)?.message ?? e).slice(0, 200);
        this.opts.logger.warn(`livebridge NATS connect failed: ${this.detail}`);
      }
      if (this.stopped) return;
      await this.sleep(backoff);
      backoff = Math.min(backoff * 2, 30_000);
    }
  }

  private sleep(ms: number): Promise<void> {
    return new Promise<void>(resolve => {
      const t = setTimeout(resolve, ms);
      t.unref?.();
    });
  }

  private subscribe(nc: NatsConnection): void {
    const subs: Subscription[] = [];
    for (const subject of SUBSCRIPTIONS) {
      subs.push(
        nc.subscribe(subject, {
          callback: (err, msg) => {
            if (err || !msg) return;
            const s = msg.subject;
            if (ignoredSubject(s)) return;
            const at = this.now();
            this.lastAt = at;
            this.meter.hit(at);
            try {
              this.opts.onMessage(s, msg.data);
            } catch (e) {
              this.opts.logger.debug(
                `livebridge message handler failed: ${(e as Error).message}`,
              );
            }
          },
        }),
      );
    }
  }

  private async watchStatus(nc: NatsConnection): Promise<void> {
    try {
      for await (const s of nc.status()) {
        switch (s.type) {
          case 'disconnect':
          case 'reconnecting':
            this.state = 'connecting';
            this.detail = `reconnecting (${s.type})`;
            break;
          case 'reconnect':
            this.state = 'connected';
            this.detail = undefined;
            break;
          case 'error':
            this.detail = String(
              (s as { error?: Error }).error?.message ?? 'error',
            ).slice(0, 200);
            break;
          default:
            break;
        }
      }
    } catch {
      // status iterator ends when the connection closes
    }
  }
}
