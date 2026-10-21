import { KeyboardEvent, useEffect } from 'react';
import { SidebarItem } from '@backstage/core-components';
import { SearchModal, useSearchModal } from '@backstage/plugin-search';
import SearchIcon from '@material-ui/icons/Search';

/** True for Cmd+K on macOS and Ctrl+K elsewhere. */
export function isPaletteShortcut(
  event: Pick<
    globalThis.KeyboardEvent,
    'key' | 'metaKey' | 'ctrlKey' | 'altKey' | 'shiftKey'
  >,
) {
  return (
    event.key.toLowerCase() === 'k' &&
    (event.metaKey || event.ctrlKey) &&
    !event.altKey &&
    !event.shiftKey
  );
}

const isMac = () =>
  typeof navigator !== 'undefined' &&
  /mac|iphone|ipad/i.test(navigator.platform);
export const paletteHint = () => (isMac() ? '⌘K' : 'Ctrl K');

const FOCUSABLE = 'input, a[href]';

/**
 * Next focus target for arrow keys inside the palette: the search box followed
 * by every result link, wrapping at both ends. Pure so it can be unit tested.
 */
export function nextIndex(current: number, count: number, key: string) {
  if (!count) return -1;
  if (key === 'ArrowDown') return current < 0 ? 0 : (current + 1) % count;
  if (key === 'ArrowUp') return current <= 0 ? count - 1 : current - 1;
  return current;
}

function onPaletteKeyDown(event: KeyboardEvent) {
  if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return;
  const target = event.target as HTMLElement;
  const dialog = target.closest('[role="dialog"]');
  if (!dialog) return;
  const stops = Array.from(dialog.querySelectorAll<HTMLElement>(FOCUSABLE));
  const next = nextIndex(stops.indexOf(target), stops.length, event.key);
  if (next >= 0) {
    event.preventDefault();
    stops[next].focus();
  }
}

/** Sidebar entry that opens the palette. Must render inside SearchModalProvider. */
export function PaletteTrigger() {
  const { toggleModal } = useSearchModal();
  return (
    <SidebarItem
      icon={SearchIcon}
      text={`Search · ${paletteHint()}`}
      onClick={toggleModal}
    />
  );
}

/**
 * Global command palette: Cmd/Ctrl-K toggles Backstage's own search modal, which
 * queries the existing search backend (catalog apps, TechDocs runbooks and docs,
 * ADRs, announcements). Arrow keys move through results; Enter opens; Esc closes.
 * Must render inside SearchModalProvider.
 */
export function CommandPalette() {
  const { state, toggleModal } = useSearchModal();
  useEffect(() => {
    const listener = (event: globalThis.KeyboardEvent) => {
      if (!isPaletteShortcut(event)) return;
      event.preventDefault();
      toggleModal();
    };
    window.addEventListener('keydown', listener);
    return () => window.removeEventListener('keydown', listener);
  }, [toggleModal]);
  return (
    <div role="presentation" onKeyDown={onPaletteKeyDown}>
      <SearchModal {...state} toggleModal={toggleModal} />
    </div>
  );
}
