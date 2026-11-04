import { render, screen } from '@testing-library/react';
import { StatusDot } from './StatusDot';

test.each([
  ['ok', 'Healthy', 'var(--onemind-status-ok)'],
  ['warn', 'Degraded', 'var(--onemind-status-warn)'],
  ['down', 'Down', 'var(--onemind-status-critical)'],
  ['unknown', 'Unknown', 'var(--onemind-status-unknown)'],
] as const)('%s is announced and drawn with its own token', (status, label, color) => {
  const { container } = render(<StatusDot status={status} showLabel />);
  expect(screen.getByRole('img', { name: `Status: ${label}` })).toBeTruthy();
  const mark = container.querySelector('span[aria-hidden]') as HTMLElement;
  expect(mark.style.background + mark.style.border).toContain(color);
});

test('unknown is hollow grey, never filled green', () => {
  const { container } = render(<StatusDot status="unknown" />);
  const mark = container.querySelector('span[aria-hidden]') as HTMLElement;
  expect(mark.style.background).toBe('transparent');
  expect(mark.style.border).toContain('--onemind-status-unknown');
  expect(mark.outerHTML).not.toContain('status-ok');
});

test('detail text is added to the accessible name', () => {
  render(<StatusDot status="down" title="platform-x: Degraded, Synced" />);
  expect(screen.getByRole('img', { name: 'Status: Down. platform-x: Degraded, Synced' })).toBeTruthy();
});
