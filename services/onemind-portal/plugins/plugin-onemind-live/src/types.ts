// Shape of GET /api/onemind/live and the status route's additions. Mirrors
// plugin-onemind-operations-backend/src/types.ts (a backend package cannot be imported here).

export type SourceStatus = 'ok' | 'stale' | 'error';
export type PanelId = 'alerts' | 'argo' | 'workloads' | 'nodes' | 'storage' | 'backups' | 'uptime' | 'events';
export type PanelState = {
  status: SourceStatus;
  observedAt: string | null;
  sources: string[];
  problems: { source: string; status: SourceStatus; error?: string }[];
};

export type LiveAlert = {
  fingerprint: string;
  name: string;
  severity: 'critical' | 'warning' | 'info';
  summary: string;
  namespace?: string;
  workload?: string;
  startsAt: string;
  runbookUrl?: string;
  grafanaUrl: string;
};
export type ArgoItem = {
  name: string;
  health: string;
  sync: string;
  reasons: string[];
  severity: 'critical' | 'warning';
  since?: string;
  message?: string;
  url: string;
};
export type WorkloadIssue = {
  namespace: string;
  pod: string;
  app?: string;
  reasons: string[];
  restarts: number;
  lastRestartAt?: string;
  severity: 'critical' | 'warning';
};
export type NodeItem = {
  name: string;
  ready: boolean;
  pressure: string[];
  cpuPercent: number | null;
  memoryPercent: number | null;
  diskPercent: number | null;
  cpuRequestedPercent: number | null;
  memoryRequestedPercent: number | null;
  pods: number;
};
export type PvcItem = {
  namespace: string;
  name: string;
  usedPercent: number;
  fillsWithin4Days: boolean;
  severity: 'warn' | 'critical';
};
export type Volsync = { sources: number; outOfSync: { namespace: string; name: string }[] };
export type BackupItem = {
  name: string;
  lastSuccess: string | null;
  suspended: boolean;
  state: 'active' | 'retired' | 'needs-setup';
  health: 'ok' | 'failing' | 'overdue' | 'never-ran' | 'suspended' | 'inactive';
  attention: 'none' | 'review' | 'setup';
  supersededBy?: string;
  pendingReason?: string;
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
export type EventItem = {
  namespace: string;
  object: string;
  reason: string;
  message: string;
  count: number;
  lastSeen: string;
};
export type AttentionItem = {
  key: string;
  kind: 'alert' | 'service' | 'workload' | 'deployment' | 'backup' | 'storage';
  severity: 'critical' | 'warning';
  title: string;
  detail: string;
  url: string;
  namespace?: string;
};
export type HealthLevel = 'ok' | 'warn' | 'critical' | 'unknown';

export type LiveResponse = {
  observedAt: string;
  cluster: string;
  health: { level: HealthLevel; headline: string; reasons: string[] };
  counts: Record<'alerts' | 'argo' | 'workloads' | 'backups' | 'backupsNeedSetup' | 'backupsRetired', number | 'unknown'>;
  panels: Record<PanelId, PanelState>;
  alerts: LiveAlert[] | null;
  argo: { items: ArgoItem[]; total: number } | null;
  workloads: { issues: WorkloadIssue[]; notReady: string[]; total: number } | null;
  nodes: NodeItem[] | null;
  storage: { pvcs: PvcItem[]; volsync: Volsync | null } | null;
  backups: { items: BackupItem[]; volsync: Volsync | null } | null;
  uptime: UptimeItem[] | null;
  events: { items: EventItem[]; windowMinutes: number } | null;
  attention: AttentionItem[];
  needsSetup: AttentionItem[];
  links: { argocd: string; headlamp: string; grafana: string };
};
