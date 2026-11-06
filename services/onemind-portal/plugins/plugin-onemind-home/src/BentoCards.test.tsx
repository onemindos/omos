import { MemoryRouter } from 'react-router-dom';
import { render, screen } from '@testing-library/react';
import { PlatformScore, scoreText } from './PlatformScore';
import { NextUpHero } from './NextUpHero';
import { TasksCard, overduePercent } from './TasksCard';

let statusResult: any;
let calendarResult: any;
let cloudResult: any;
jest.mock('./status', () => ({
  ...jest.requireActual('./status'),
  useOperationsStatus: () => statusResult,
}));
jest.mock('./CalendarPanel', () => ({
  useCalendarSummary: () => calendarResult,
}));
jest.mock('./cloudWork', () => ({ useCloudWork: () => cloudResult }));

const wrap = (node: any) => render(<MemoryRouter>{node}</MemoryRouter>);

test('score shows an em dash and "unavailable" without data', () => {
  statusResult = { data: undefined, error: true };
  wrap(<PlatformScore />);
  expect(scoreText(undefined)).toBe('—');
  expect(screen.getByText('unavailable')).not.toBeNull();
  expect(screen.queryByText('100')).toBeNull();
});

test('score shows the number and the top attention item', () => {
  statusResult = {
    error: false,
    data: {
      observedAt: new Date().toISOString(),
      partial: false,
      totalWorkloads: 1,
      unavailable: [],
      deployments: [],
      backups: [],
      attention: [
        {
          key: 'a',
          kind: 'x',
          severity: 'warning',
          title: 'Disk low',
          detail: '',
          url: '/',
        },
      ],
    },
  };
  wrap(<PlatformScore />);
  expect(screen.getByText('92')).not.toBeNull();
  expect(screen.getByText(/Disk low/)).not.toBeNull();
});

test('next up shows the empty state and no timeline blocks when nothing is scheduled', () => {
  calendarResult = {
    error: '',
    data: {
      state: 'ok',
      coverage: '',
      observedAt: null,
      expiresAt: null,
      events: [],
      link: 'https://calendar.zoho.com/',
    },
  };
  cloudResult = {
    error: false,
    data: {
      state: 'ok',
      observedAt: null,
      expiresAt: null,
      tasks: [],
      link: 'https://projects.zoho.com',
    },
  };
  wrap(<NextUpHero />);
  expect(screen.getByText(/Nothing scheduled/)).not.toBeNull();
  expect(screen.getByRole('img').getAttribute('aria-label')).toContain(
    'no remaining events',
  );
});

test('next up draws no timeline when the calendar is not configured', () => {
  calendarResult = {
    error: '',
    data: {
      state: 'configuration-required',
      coverage: '',
      observedAt: null,
      expiresAt: null,
      events: [],
      link: '',
    },
  };
  cloudResult = { error: false, data: undefined };
  wrap(<NextUpHero />);
  expect(screen.queryByRole('img')).toBeNull();
  expect(screen.getByText(/Connect read-only Calendar access/)).not.toBeNull();
});

test('tasks show counts, a proportion bar and honest unavailable state', () => {
  cloudResult = {
    error: false,
    data: {
      state: 'ok',
      observedAt: null,
      expiresAt: null,
      open: 8,
      overdue: 2,
      tasks: [
        { id: '1', title: 'Ship it', dueDate: '2020-01-01', overdue: true },
      ],
      link: '',
    },
  };
  wrap(<TasksCard />);
  expect(screen.getByRole('img').getAttribute('aria-label')).toContain(
    '2 of 8',
  );
  expect(screen.getByText('Ship it')).not.toBeNull();
  expect(overduePercent(8, 2)).toBe(25);
  expect(overduePercent(undefined, 2)).toBeUndefined();
});

test('tasks never invent counts when the provider is down', () => {
  cloudResult = { error: true, data: undefined };
  wrap(<TasksCard />);
  expect(screen.queryByRole('img')).toBeNull();
  expect(screen.getByText(/summary unavailable/)).not.toBeNull();
});
