import { fireEvent, render, screen } from '@testing-library/react';
import { GitHubGate, githubGateModules, resetGitHubConnection } from './githubGate';

beforeEach(() => resetGitHubConnection());

describe('GitHubGate', () => {
  it('does not render the plugin components, so no sign-in is requested, until Connect GitHub is clicked', () => {
    const mounted = jest.fn();
    function Child() {
      mounted();
      return <div>pull requests</div>;
    }
    render(
      <GitHubGate>
        <Child />
      </GitHubGate>,
    );
    expect(screen.getByText(/GitHub is not connected/)).toBeTruthy();
    expect(screen.queryByText('pull requests')).toBeNull();
    expect(mounted).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole('button', { name: 'Connect GitHub' }));
    expect(screen.getByText('pull requests')).toBeTruthy();
    expect(mounted).toHaveBeenCalled();
  });

  it('one click connects every gate on the page, and the choice is remembered', () => {
    const { unmount } = render(
      <>
        <GitHubGate>
          <span>card one</span>
        </GitHubGate>
        <GitHubGate>
          <span>card two</span>
        </GitHubGate>
      </>,
    );
    fireEvent.click(screen.getAllByRole('button', { name: 'Connect GitHub' })[0]);
    expect(screen.getByText('card one')).toBeTruthy();
    expect(screen.getByText('card two')).toBeTruthy();
    unmount();

    render(
      <GitHubGate>
        <span>later visit</span>
      </GitHubGate>,
    );
    expect(screen.getByText('later visit')).toBeTruthy();
  });

  it('keeps working when browser storage is blocked', () => {
    const original = Object.getOwnPropertyDescriptor(window, 'localStorage');
    Object.defineProperty(window, 'localStorage', {
      configurable: true,
      get() {
        throw new Error('blocked');
      },
    });
    try {
      render(
        <GitHubGate>
          <span>works</span>
        </GitHubGate>,
      );
      expect(screen.queryByText('works')).toBeNull();
      fireEvent.click(screen.getByRole('button', { name: 'Connect GitHub' }));
      expect(screen.getByText('works')).toBeTruthy();
    } finally {
      if (original) Object.defineProperty(window, 'localStorage', original);
    }
  });

  it('wraps both GitHub plugins', () => {
    expect(githubGateModules).toHaveLength(2);
  });
});
