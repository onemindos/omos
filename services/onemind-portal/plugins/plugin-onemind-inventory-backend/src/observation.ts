/** The latest successful cluster observation, shared by the provider and the processor of one replica. */
export class ObservationStore {
  private current?: { at: string; keys: Set<string> };

  set(at: string, keys: Iterable<string>): void {
    this.current = { at, keys: new Set(keys) };
  }

  /** The observation when it is younger than maxAgeMs, otherwise undefined (never guess from stale data). */
  fresh(
    now: Date,
    maxAgeMs: number,
  ): { at: string; keys: Set<string> } | undefined {
    if (!this.current) return undefined;
    const age = now.getTime() - Date.parse(this.current.at);
    return Number.isFinite(age) && age <= maxAgeMs ? this.current : undefined;
  }
}
