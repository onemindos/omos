import { createElement } from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { MockEntityListContextProvider } from '@backstage/plugin-catalog-react/testUtils';
import {
  EntityPillarFilter,
  EntityRoleFilter,
  PILLAR_LABEL,
  PILLAR_REVIEW,
  PillarPicker,
  ROLE_LABEL,
  ViewPicker,
  currentView,
  pillarOptions,
  viewFilters,
  viewOptions,
} from './catalogViewFilters';

describe('pillar filter contract', () => {
  it.each(['so', 'le', 'hp', 'ge'] as const)(
    '%s matches the reviewed primary label only',
    pillar => {
      expect(new EntityPillarFilter(pillar).getCatalogFilters()).toEqual({
        [PILLAR_LABEL]: pillar,
      });
    },
  );
  it('Shared is the reviewed shared label; Unclassified is a positive state, never a guess', () => {
    expect(new EntityPillarFilter('shared').getCatalogFilters()).toEqual({
      [PILLAR_LABEL]: 'shared',
    });
    expect(new EntityPillarFilter('unclassified').getCatalogFilters()).toEqual({
      [PILLAR_REVIEW]: 'unclassified',
    });
    expect(new EntityPillarFilter('all').getCatalogFilters()).toEqual({});
  });
  it('offers All, SO, LE, HP, GE, Shared and Unclassified in that order', () => {
    expect(pillarOptions.map(o => o.id)).toEqual([
      'all',
      'so',
      'le',
      'hp',
      'ge',
      'shared',
      'unclassified',
    ]);
  });
});

describe('saved views', () => {
  it('map to exact kind, type and role filters', () => {
    expect(viewFilters('apps').kind?.getCatalogFilters()).toEqual({
      kind: 'component',
    });
    expect(viewFilters('apps').role?.getCatalogFilters()).toEqual({
      [ROLE_LABEL]: 'app',
    });
    expect(viewFilters('infrastructure').role?.value).toBe('workload');
    const repos = viewFilters('repos');
    expect(repos.kind?.value).toBe('resource');
    expect(repos.type?.getCatalogFilters()).toEqual({
      'spec.type': ['code-repository'],
    });
    expect(viewFilters('everything').role).toBeUndefined();
    expect(viewOptions.map(o => o.label)).toEqual([
      'Apps',
      'Infrastructure',
      'Repos',
      'Everything',
    ]);
  });
  it('currentView round-trips every preset', () => {
    for (const { id } of viewOptions)
      expect(currentView(viewFilters(id) as never)).toBe(id);
  });
});

describe('pickers', () => {
  it('Pillar picker applies and clears the filter', () => {
    const updateFilters = jest.fn();
    render(
      createElement(
        MockEntityListContextProvider,
        { value: { filters: {}, updateFilters } },
        createElement(PillarPicker),
      ),
    );
    fireEvent.click(screen.getByText('Unclassified'));
    expect(updateFilters).toHaveBeenLastCalledWith({
      pillar: expect.any(EntityPillarFilter),
    });
    expect(updateFilters.mock.calls[0][0].pillar.value).toBe('unclassified');
    fireEvent.click(screen.getByText('All'));
    expect(updateFilters).toHaveBeenLastCalledWith({ pillar: undefined });
  });

  it('View picker starts on Apps, and a deep link to another list is left alone', () => {
    const fresh = jest.fn();
    render(
      createElement(
        MockEntityListContextProvider,
        { value: { filters: {}, updateFilters: fresh } },
        createElement(ViewPicker),
      ),
    );
    expect(fresh).toHaveBeenCalledTimes(1);
    expect(fresh.mock.calls[0][0].role).toBeInstanceOf(EntityRoleFilter);
    const linked = jest.fn();
    render(
      createElement(
        MockEntityListContextProvider,
        {
          value: {
            filters: { kind: viewFilters('repos').kind },
            updateFilters: linked,
          },
        },
        createElement(ViewPicker),
      ),
    );
    expect(linked).not.toHaveBeenCalled();
  });
});

// Runs the System-redirect routing tests inside the step the gate already names.
import './systemRedirects.test';
