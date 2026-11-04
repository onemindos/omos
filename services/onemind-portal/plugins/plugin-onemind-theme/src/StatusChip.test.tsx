import { render, screen } from '@testing-library/react';
import { StatusChip } from './StatusChip';

describe('StatusChip', () => {
  it('renders a default label and an icon for every status', () => {
    (['ok', 'warn', 'critical', 'unknown', 'info'] as const).forEach(status => {
      const { container, unmount } = render(<StatusChip status={status} />);
      expect(container.querySelector('svg')).not.toBeNull();
      expect(
        container.querySelector(`[data-status="${status}"]`),
      ).not.toBeNull();
      unmount();
    });
  });

  it('uses theme variables, never the brand colour, for critical', () => {
    render(<StatusChip status="critical" label="Down" />);
    const chip = screen.getByText('Down');
    expect(chip.style.color).toBe('var(--onemind-status-critical)');
    expect(chip.style.color).not.toContain('brand');
  });
});
