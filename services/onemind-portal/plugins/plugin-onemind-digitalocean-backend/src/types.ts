/** What the summary route returns. Nothing here is, or is derived from, the API token. */
export type SectionState = 'ok' | 'unavailable' | 'not-available';

export interface Section<T> {
  state: SectionState;
  /** Why a section is not 'ok', in plain words (never a token, header or raw upstream body). */
  reason?: string;
  items: T[];
}

export interface PvcLink {
  namespace: string;
  name: string;
}
export interface AppLink {
  /** Catalog entity ref of the owning component, e.g. component:default/k8s-... */
  ref: string;
  title?: string;
}

export interface NodePool {
  id: string;
  name: string;
  size: string;
  count: number;
  autoScale: boolean;
  minNodes?: number;
  maxNodes?: number;
}
export interface Cluster {
  id: string;
  name: string;
  region: string;
  version: string;
  status?: string;
  highAvailability: boolean;
  /** Newest version DOKS offers this cluster, or null when it is up to date. */
  upgradeAvailable: string | null;
  upgradeCheck: 'ok' | 'unavailable';
  createdAt?: string;
  nodePools: NodePool[];
}
export interface Droplet {
  id: string;
  name: string;
  status: string;
  size: string;
  region: string;
  vcpus?: number;
  memoryMb?: number;
  diskGb?: number;
  tags: string[];
  createdAt?: string;
}
export interface Volume {
  id: string;
  name: string;
  sizeGb: number;
  region: string;
  attachedDropletIds: string[];
  /** True when no droplet has it attached: a candidate for review (report only). */
  unattached: boolean;
  createdAt?: string;
  pvc?: PvcLink;
  app?: AppLink;
}
export interface Snapshot {
  id: string;
  name: string;
  resourceType: string;
  resourceId?: string;
  sizeGb: number;
  minDiskGb?: number;
  regions: string[];
  createdAt?: string;
  ageDays: number | null;
}
export interface SnapshotStats {
  count: number;
  totalSizeGb: number;
  oldest: Snapshot | null;
  largest: Snapshot | null;
  /** Snapshots older than this many days, for the review list. */
  olderThanDays: number;
  olderCount: number;
}
export interface LoadBalancer {
  id: string;
  name: string;
  region: string;
  size?: string;
  type?: string;
  status: string;
  forwardingRules: number;
  dropletCount: number;
  /** No droplets and no tag attached: an idle balancer (report only). */
  idle: boolean;
  createdAt?: string;
}
export interface Firewall {
  id: string;
  name: string;
  status: string;
  inboundRules: number;
  outboundRules: number;
  dropletCount: number;
  tags: string[];
}
export interface Vpc {
  id: string;
  name: string;
  region: string;
  ipRange: string;
  isDefault: boolean;
}
export interface Domain {
  name: string;
  /** null when the record list could not be read. */
  recordCount: number | null;
}
export interface UptimeCheck {
  id: string;
  name: string;
  type: string;
  target: string;
  enabled: boolean;
  regions: string[];
  state: 'up' | 'down' | 'unknown';
  regionStates: Record<string, string>;
}
export interface AlertPolicy {
  id: string;
  description: string;
  type: string;
  enabled: boolean;
  compare?: string;
  value?: number;
  window?: string;
}
export interface Invoice {
  period: string;
  amount: string;
}
export interface Billing {
  state: SectionState;
  reason?: string;
  /** month_to_date_usage: what was used this month BEFORE credits. */
  usageBeforeCredits?: string;
  /** month_to_date_balance: what is owed after credits. */
  amountDue?: string;
  accountBalance?: string;
  /** usageBeforeCredits scaled to the whole month, in USD with two decimals. */
  projectedMonth?: string;
  daysElapsed?: number;
  daysInMonth?: number;
  generatedAt?: string;
  invoices: Invoice[];
  /** The open invoice for the current month, when the API offers it. */
  invoicePreview?: Invoice;
}
export interface AccountInfo {
  state: SectionState;
  reason?: string;
  status?: string;
  teamName?: string;
  dropletLimit?: number;
  volumeLimit?: number;
}

export type ResourceType =
  | 'droplet'
  | 'volume'
  | 'snapshot'
  | 'load-balancer'
  | 'firewall'
  | 'vpc'
  | 'kubernetes-cluster'
  | 'node-pool';

export interface Change {
  type: ResourceType;
  id: string;
  name: string;
}
export interface ChangeSet {
  at: string;
  added: Change[];
  removed: Change[];
}

export interface Snapshots {
  generatedAt: string;
  account: AccountInfo;
  billing: Billing;
  kubernetes: Section<Cluster>;
  droplets: Section<Droplet>;
  volumes: Section<Volume> & {
    unattachedCount: number;
    unattachedSizeGb: number;
  };
  snapshots: Section<Snapshot> & { stats: SnapshotStats };
  loadBalancers: Section<LoadBalancer>;
  firewalls: Section<Firewall>;
  vpcs: Section<Vpc>;
  domains: Section<Domain>;
  uptime: Section<UptimeCheck>;
  alerts: Section<AlertPolicy>;
  /** The DO API v2 cannot list Spaces buckets with an API token. */
  spaces: Section<never>;
}

export interface Summary extends Snapshots {
  state: 'not-connected' | 'connected';
  /** When the cached data is next refreshed. */
  expiresAt: string | null;
  /** Newest first, kept in memory since the portal started. */
  changes: ChangeSet[];
}
