/** Messages per second over a sliding 10 s window, in 1 s buckets. */
export class RateMeter {
  private readonly counts = new Array<number>(10).fill(0);
  private readonly seconds = new Array<number>(10).fill(-1);

  hit(nowMs: number, n = 1): void {
    const s = Math.floor(nowMs / 1000);
    const i = s % 10;
    if (this.seconds[i] !== s) {
      this.seconds[i] = s;
      this.counts[i] = 0;
    }
    this.counts[i] += n;
  }

  rate(nowMs: number): number {
    const s = Math.floor(nowMs / 1000);
    let total = 0;
    for (let i = 0; i < 10; i++) {
      if (this.seconds[i] > s - 10 && this.seconds[i] <= s)
        total += this.counts[i];
    }
    return Math.round((total / 10) * 10) / 10;
  }
}
