import type { BackupItem } from './backupState';
import type { NodeUsage, PvcItem, VolsyncSummary } from './liveQueries';
import type { SourceId, SourceStatus } from './source';

export type PanelId =
  | 'alerts'
  | 'argo'
  | 'workloads'
  | 'nodes'
  | 'storage'
  | 'backups'
  | 'uptime'
  | 'events';

export type PanelState = {
  status: SourceStatus;
  observedAt: string | null;
  sources: SourceId[];
  /** Short reason for each source that is not ok. */
  problems: { source: SourceId; status: SourceStatus; error?: string }[];
};

export type AlertSeverity = 'critical' | 'warning' | 'info';

export type LiveAlert = {
  fingerprint: string;
  name: string;
  severity: AlertSeverity;
  summary: string;
  namespace?: string;
  /** Best guess at the owning workload from the alert labels. */
  workload?: string;
  startsAt: string;
  runbookUrl?: string;
  grafanaUrl: string;
};

export type ArgoReason = 'Degraded' | 'OutOfSync' | 'Progressing' | 'Sync failed' | 'Unknown';
export type ArgoItem = {
  name: string;
  health: string;
  sync: string;
  reasons: ArgoReason[];
  severity: 'critical' | 'warning';
  /** ISO time the application entered its current health state, when Argo CD reports it. */
  since?: string;
  message?: string;
  url: string;
};

export type WorkloadReason =
  | 'CrashLoopBackOff'
  | 'ImagePullBackOff'
  | 'Pending'
  | 'NotReady'
  | 'OOMKilled'
  | 'Restarted'
  | 'ConfigError';
export type WorkloadIssue = {
  namespace: string;
  pod: string;
  /** `app.kubernetes.io/instance` or `app` label when present. */
  app?: string;
  reasons: WorkloadReason[];
  restarts: number;
  lastRestartAt?: string;
  severity: 'critical' | 'warning';
};

export type NodeItem = NodeUsage & {
  name: string;
  ready: boolean;
  pressure: string[];
  cpuRequestedPercent: number | null;
  memoryRequestedPercent: number | null;
  pods: number;
};

export type EventItem = {
  namespace: string;
  object: string;
  reason: string;
  message: string;
  count: number;
  lastSeen: string;
};

export type UptimeItem = {
  name: string;
  state: string;
  scope: string;
  link: string;
  entity: string;
  httpStatus?: number;
  latencyMs?: number;
  observedAt: string;
};

export type AttentionKind = 'alert' | 'service' | 'workload' | 'deployment' | 'backup' | 'storage';
export type AttentionItem = {
  key: string;
  kind: AttentionKind;
  severity: 'critical' | 'warning';
  title: string;
  detail: string;
  url: string;
  namespace?: string;
};

export type HealthLevel = 'ok' | 'warn' | 'critical' | 'unknown';
export type HealthBanner = { level: HealthLevel; headline: string; reasons: string[] };

export type StripCount = number | 'unknown';
export type StatusCounts = {
  alerts: StripCount;
  argo: StripCount;
  workloads: StripCount;
  backups: StripCount;
  backupsNeedSetup: StripCount;
  backupsRetired: StripCount;
};

export type LiveResponse = {
  observedAt: string;
  cluster: string;
  health: HealthBanner;
  counts: StatusCounts;
  panels: Record<PanelId, PanelState>;
  alerts: LiveAlert[] | null;
  argo: { items: ArgoItem[]; total: number } | null;
  workloads: { issues: WorkloadIssue[]; notReady: string[]; total: number } | null;
  nodes: NodeItem[] | null;
  storage: { pvcs: PvcItem[]; volsync: VolsyncSummary | null } | null;
  backups: { items: BackupItem[]; volsync: VolsyncSummary | null } | null;
  uptime: UptimeItem[] | null;
  events: { items: EventItem[]; windowMinutes: number } | null;
  /** Items that need action now. Home status strip, Today and this page all read this list. */
  attention: AttentionItem[];
  /** Backups waiting for setup only a person can do. Separate from attention. */
  needsSetup: AttentionItem[];
  links: { argocd: string; headlamp: string; grafana: string };
};
