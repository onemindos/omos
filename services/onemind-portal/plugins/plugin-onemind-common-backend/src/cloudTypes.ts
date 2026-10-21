/** Server-to-admin contract. Never include raw provider objects, tokens or descriptions. */
export type CloudTask = {
  id: string;
  title: string;
  dueDate: string | null;
  overdue: boolean;
};
export type CloudSummary = {
  providerId: string;
  state: 'configuration-required' | 'current' | 'stale' | 'unavailable';
  source: string;
  observedAt: string | null;
  expiresAt: string | null;
  attemptedAt: string;
  open?: number;
  overdue?: number;
  sampled?: number;
  pages?: number;
  partial?: boolean;
  coverage: string;
  tasks: CloudTask[];
  link: string;
};
/** Future providers must own their read-only credentials, sanitization and cache. */
export interface CloudProvider {
  readonly id: string;
  summary(): Promise<CloudSummary>;
}
