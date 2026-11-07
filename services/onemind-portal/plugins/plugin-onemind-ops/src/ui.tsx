import type { ReactNode } from 'react';
import { makeStyles, useTheme } from '@material-ui/core/styles';
import { categoryPalette } from '@internal/plugin-onemind-theme';
import type { CategoryColor } from '@internal/plugin-onemind-theme';
import type { CellHealth, ConnectionState } from './contract';

/** Theme variables from `@internal/plugin-onemind-theme` (brand.css). */
export const cssVar = (name: string) => `var(--onemind-${name})`;

export const KIND_ORDER = [
  'tel',
  'det',
  'evt',
  'cmd',
  'intel',
  'ent',
  'tak',
  'agents',
] as const;

const KIND_PALETTE: Record<string, CategoryColor> = {
  tel: categoryPalette.cyan,
  det: categoryPalette.violet,
  evt: categoryPalette.amber,
  cmd: categoryPalette.orange,
  intel: categoryPalette.teal,
  ent: categoryPalette.indigo,
  tak: categoryPalette.moss,
  agents: categoryPalette.plum,
};

/** Identity colour of a bus subject kind. Unknown kinds get slate. */
export function useKindColor(): (kind: string) => string {
  const dark = useTheme().palette.type === 'dark';
  return kind => {
    const c = KIND_PALETTE[kind] ?? categoryPalette.slate;
    return dark ? c.bright : c.base;
  };
}

export type StatusVar = 'ok' | 'warn' | 'critical' | 'unknown' | 'info';

export const CONNECTION_STATUS: Record<ConnectionState, StatusVar> = {
  connected: 'ok',
  connecting: 'warn',
  down: 'critical',
  disabled: 'unknown',
};

export const CELL_STATUS: Record<CellHealth, StatusVar> = {
  ok: 'ok',
  warn: 'warn',
  down: 'critical',
};

export const statusFg = (s: StatusVar) => cssVar(`status-${s}`);
export const statusBg = (s: StatusVar) => cssVar(`status-${s}-bg`);

const useStyles = makeStyles({
  panel: {
    display: 'flex',
    flexDirection: 'column',
    minHeight: 0,
    background: cssVar('surface-1'),
    border: `1px solid ${cssVar('border')}`,
    borderRadius: 6,
    overflow: 'hidden',
  },
  head: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    padding: '6px 10px',
    borderBottom: `1px solid ${cssVar('border')}`,
    background: cssVar('surface-2'),
    fontSize: 11,
    fontWeight: 700,
    letterSpacing: '0.08em',
    textTransform: 'uppercase',
    color: cssVar('text-2'),
  },
  grow: { flex: 1 },
  body: { flex: 1, minHeight: 0, overflow: 'auto' },
  note: { padding: 12, fontSize: 13, color: cssVar('text-2') },
});

export function Panel(props: {
  title: string;
  aside?: ReactNode;
  children: ReactNode;
  className?: string;
  bodyStyle?: React.CSSProperties;
  label?: string;
}) {
  const c = useStyles();
  return (
    <section
      className={`${c.panel} ${props.className ?? ''}`}
      aria-label={props.label ?? props.title}
    >
      <header className={c.head}>
        <span>{props.title}</span>
        <span className={c.grow} />
        {props.aside}
      </header>
      <div className={c.body} style={props.bodyStyle}>
        {props.children}
      </div>
    </section>
  );
}

/** Message used for an empty or unavailable panel. Never a bare zero. */
export function Note({
  children,
  tone,
}: {
  children: ReactNode;
  tone?: StatusVar;
}) {
  const c = useStyles();
  return (
    <div
      className={c.note}
      role="status"
      style={tone ? { color: statusFg(tone) } : undefined}
    >
      {children}
    </div>
  );
}

export function Dot({
  status,
  size = 9,
  pulse,
}: {
  status: StatusVar;
  size?: number;
  pulse?: boolean;
}) {
  return (
    <span
      aria-hidden
      data-status={status}
      style={{
        display: 'inline-block',
        width: size,
        height: size,
        borderRadius: '50%',
        background: statusFg(status),
        boxShadow: pulse ? `0 0 6px ${statusFg(status)}` : undefined,
        flex: 'none',
      }}
    />
  );
}

/** Tiny inline line chart. */
export function Sparkline({
  values,
  color,
  width = 64,
  height = 16,
}: {
  values: number[];
  color: string;
  width?: number;
  height?: number;
}) {
  const max = Math.max(1, ...values);
  const step = values.length > 1 ? width / (values.length - 1) : width;
  const pts = values
    .map(
      (v, i) =>
        `${(i * step).toFixed(1)},${(
          height -
          1 -
          (v / max) * (height - 2)
        ).toFixed(1)}`,
    )
    .join(' ');
  return (
    <svg
      width={width}
      height={height}
      viewBox={`0 0 ${width} ${height}`}
      aria-hidden
      style={{ display: 'block' }}
    >
      <polyline
        points={pts}
        fill="none"
        stroke={color}
        strokeWidth={1.2}
        strokeLinejoin="round"
      />
    </svg>
  );
}
