import { act, render } from '@testing-library/react';
import { useAlertRefresh, ALERT_SIGNAL_CHANNEL } from './useAlertRefresh';

let signal: { lastSignal: unknown; isSignalsAvailable: boolean } = { lastSignal: null, isSignalsAvailable: true };
const channels: string[] = [];
jest.mock('@backstage/plugin-signals-react', () => ({
  useSignal: (channel: string) => {
    channels.push(channel);
    return signal;
  },
}));

function Probe({ refresh }: { refresh: () => void }) {
  const { push } = useAlertRefresh(refresh);
  return <span>{push ? 'push' : 'poll'}</span>;
}

beforeEach(() => {
  jest.useFakeTimers();
  signal = { lastSignal: null, isSignalsAvailable: true };
  channels.length = 0;
});
afterEach(() => jest.useRealTimers());

it('subscribes to the F3 alerts channel', () => {
  render(<Probe refresh={() => {}} />);
  expect(ALERT_SIGNAL_CHANNEL).toBe('onemind:alerts');
  expect(channels).toContain('onemind:alerts');
});

it('refreshes on mount and every 30 seconds without Signals', () => {
  signal = { lastSignal: null, isSignalsAvailable: false };
  const refresh = jest.fn();
  const view = render(<Probe refresh={refresh} />);
  expect(view.container.textContent).toBe('poll');
  expect(refresh).toHaveBeenCalledTimes(1);
  act(() => {
    jest.advanceTimersByTime(30000);
  });
  expect(refresh).toHaveBeenCalledTimes(2);
  act(() => {
    jest.advanceTimersByTime(60000);
  });
  expect(refresh).toHaveBeenCalledTimes(4);
});

it('refreshes right after an alert Signal, once per burst', () => {
  const refresh = jest.fn();
  const view = render(<Probe refresh={refresh} />);
  expect(view.container.textContent).toBe('push');
  refresh.mockClear();
  signal = { lastSignal: { action: 'alert.firing', key: 'a' }, isSignalsAvailable: true };
  view.rerender(<Probe refresh={refresh} />);
  signal = { lastSignal: { action: 'alert.firing', key: 'b' }, isSignalsAvailable: true };
  view.rerender(<Probe refresh={refresh} />);
  act(() => {
    jest.advanceTimersByTime(500);
  });
  expect(refresh).toHaveBeenCalledTimes(1);
});

it('stops polling when unmounted', () => {
  const refresh = jest.fn();
  const view = render(<Probe refresh={refresh} />);
  view.unmount();
  refresh.mockClear();
  act(() => {
    jest.advanceTimersByTime(120000);
  });
  expect(refresh).not.toHaveBeenCalled();
});
