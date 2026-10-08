import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { ActiveStack } from './ActiveStack';

const set = jest.fn();
let saved: Record<string, unknown> = {};
const bucket = {
  snapshot: () => ({ presence: 'present', value: saved }),
  observe$: () => ({ subscribe: () => ({ unsubscribe: jest.fn() }) }),
  set,
};
jest.mock('@backstage/core-plugin-api', () => ({
  storageApiRef: {},
  useApi: () => ({ forBucket: () => bucket }),
}));
jest.mock('./useDirectory', () => ({
  useDirectory: () => ({
    loaded: true,
    error: false,
    all: [
      {
        ref: 'component:default/notes',
        title: 'Test Notes',
        url: 'https://example.com',
        groupTitle: 'Knowledge',
      },
    ],
  }),
}));

beforeEach(() => {
  saved = {};
  set.mockReset();
  set.mockResolvedValue(undefined);
});

test('starts unassigned and saves a purpose decision without replacing other purposes', async () => {
  saved = {
    Email: { apps: [], data: 'Mail archive', migration: '', history: [] },
  };
  render(<ActiveStack />);
  expect(screen.getAllByText('Not assigned')).toHaveLength(11);
  fireEvent.click(screen.getByRole('button', { name: 'Manage Notes' }));
  fireEvent.click(screen.getByRole('button', { name: 'Add app' }));
  fireEvent.mouseDown(screen.getByLabelText('App'));
  fireEvent.click(
    await screen.findByRole('option', { name: 'Test Notes · Knowledge' }),
  );
  fireEvent.change(screen.getByLabelText('What I use this app for'), {
    target: { value: 'Daily writing' },
  });
  fireEvent.change(screen.getByLabelText('Reason for this decision'), {
    target: { value: 'My daily tool' },
  });
  fireEvent.click(screen.getByRole('button', { name: 'Save decision' }));
  await waitFor(() => expect(set).toHaveBeenCalled());
  const value = set.mock.calls[0][1];
  expect(value.Email.data).toBe('Mail archive');
  expect(value.Notes.apps[0]).toEqual({
    ref: 'component:default/notes',
    title: 'Test Notes',
    role: 'Primary',
    use: 'Daily writing',
  });
  expect(value.Notes.history[0].reason).toBe('My daily tool');
});

test('preserves a saved choice whose catalog app disappears', () => {
  saved = {
    Notes: {
      apps: [
        { ref: 'component:default/missing', role: 'Primary', use: 'Research' },
      ],
      data: '',
      migration: '',
      history: [],
    },
  };
  render(<ActiveStack />);
  expect(screen.getByText('App no longer in directory')).toBeTruthy();
  expect(screen.getByText('Research')).toBeTruthy();
});

test('records a native app without requiring a web link and reports failed saves', async () => {
  set.mockRejectedValue(new Error('offline'));
  render(<ActiveStack />);
  fireEvent.click(screen.getByRole('button', { name: 'Manage Notes' }));
  fireEvent.click(screen.getByRole('button', { name: 'Add app' }));
  fireEvent.mouseDown(screen.getByLabelText('App'));
  fireEvent.click(
    await screen.findByRole('option', {
      name: 'Native or other app — enter a name',
    }),
  );
  fireEvent.change(screen.getByLabelText('App name'), {
    target: { value: 'Apple Notes' },
  });
  fireEvent.change(screen.getByLabelText('What I use this app for'), {
    target: { value: 'Quick capture' },
  });
  fireEvent.change(screen.getByLabelText('Reason for this decision'), {
    target: { value: 'Current tool' },
  });
  fireEvent.click(screen.getByRole('button', { name: 'Save decision' }));
  expect(await screen.findByRole('alert')).toHaveProperty(
    'textContent',
    'Your decision could not be saved. Try again.',
  );
  expect(set.mock.calls[0][1].Notes.apps[0].title).toBe('Apple Notes');
  expect(screen.getByRole('button', { name: 'Save decision' })).toBeTruthy();
});
