import { fireEvent, render, screen } from '@testing-library/react';
import { CommandPalette, isPaletteShortcut, nextIndex } from './CommandPalette';

const toggleModal = jest.fn();
jest.mock('@backstage/plugin-search', () => ({
  useSearchModal: () => ({ state: { open: true, hidden: false }, toggleModal }),
  SearchModal: ({ toggleModal: toggle }: { toggleModal: () => void }) => (
    <div role="dialog" aria-label="Search Modal">
      <input aria-label="search" />
      <a href="/one">one</a>
      <a href="/two">two</a>
      <button onClick={toggle}>close</button>
    </div>
  ),
}));
jest.mock('@backstage/core-components', () => ({ SidebarItem: () => null }));

beforeEach(() => toggleModal.mockReset());

describe('isPaletteShortcut', () => {
  const base = {
    key: 'k',
    metaKey: false,
    ctrlKey: false,
    altKey: false,
    shiftKey: false,
  };
  it('accepts Cmd-K and Ctrl-K', () => {
    expect(isPaletteShortcut({ ...base, metaKey: true })).toBe(true);
    expect(isPaletteShortcut({ ...base, ctrlKey: true, key: 'K' })).toBe(true);
  });
  it('ignores plain k and other chords', () => {
    expect(isPaletteShortcut(base)).toBe(false);
    expect(isPaletteShortcut({ ...base, metaKey: true, shiftKey: true })).toBe(
      false,
    );
    expect(isPaletteShortcut({ ...base, ctrlKey: true, key: 'j' })).toBe(false);
  });
});

describe('nextIndex', () => {
  it('moves down and up with wrap-around', () => {
    expect(nextIndex(-1, 3, 'ArrowDown')).toBe(0);
    expect(nextIndex(2, 3, 'ArrowDown')).toBe(0);
    expect(nextIndex(0, 3, 'ArrowUp')).toBe(2);
    expect(nextIndex(2, 3, 'ArrowUp')).toBe(1);
    expect(nextIndex(0, 0, 'ArrowDown')).toBe(-1);
  });
});

describe('CommandPalette', () => {
  it('toggles on Ctrl-K and Cmd-K anywhere on the page', () => {
    render(<CommandPalette />);
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true });
    fireEvent.keyDown(window, { key: 'k', metaKey: true });
    fireEvent.keyDown(window, { key: 'k' });
    expect(toggleModal).toHaveBeenCalledTimes(2);
  });

  it('moves focus through the search box and results with the arrow keys', () => {
    render(<CommandPalette />);
    const input = screen.getByLabelText('search');
    input.focus();
    fireEvent.keyDown(input, { key: 'ArrowDown' });
    expect(document.activeElement).toBe(screen.getByText('one'));
    fireEvent.keyDown(document.activeElement!, { key: 'ArrowDown' });
    expect(document.activeElement).toBe(screen.getByText('two'));
    fireEvent.keyDown(document.activeElement!, { key: 'ArrowDown' });
    expect(document.activeElement).toBe(input);
    fireEvent.keyDown(input, { key: 'ArrowUp' });
    expect(document.activeElement).toBe(screen.getByText('two'));
  });
});
