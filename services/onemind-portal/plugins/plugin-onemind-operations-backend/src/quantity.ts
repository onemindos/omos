// Kubernetes resource quantities, reduced to the numbers the Live page needs.
const SUFFIX: Record<string, number> = {
  n: 1e-9,
  u: 1e-6,
  m: 1e-3,
  '': 1,
  k: 1e3,
  K: 1e3,
  M: 1e6,
  G: 1e9,
  T: 1e12,
  P: 1e15,
  E: 1e18,
  Ki: 1024,
  Mi: 1024 ** 2,
  Gi: 1024 ** 3,
  Ti: 1024 ** 4,
  Pi: 1024 ** 5,
  Ei: 1024 ** 6,
};

/** Parses `250m`, `2`, `16Gi`, `1e3`. Returns undefined for anything unrecognized. */
export function parseQuantity(value: string | number | undefined) {
  if (typeof value === 'number') return Number.isFinite(value) ? value : undefined;
  const match = /^([+-]?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)([A-Za-z]{0,2})$/.exec(
    value?.trim() ?? '',
  );
  if (!match) return undefined;
  const unit = SUFFIX[match[2]];
  return unit === undefined ? undefined : Number(match[1]) * unit;
}
