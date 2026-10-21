import { render, screen } from '@testing-library/react';
import { configApiRef } from '@backstage/core-plugin-api';
import { mockApis, TestApiProvider } from '@backstage/frontend-test-utils';
import ComponentOverview from './ComponentOverview';

const workflowMounted = jest.fn();
jest.mock('@openchoreo/backstage-plugin', () => ({
  DeploymentStatusCard: () => <div>Deployment status</div>,
  EntityWarningStrip: () => <div>Entity warnings</div>,
  RuntimeHealthCard: () => <div>Runtime health</div>,
  WorkflowsOverviewCard: () => {
    workflowMounted();
    return <div>Workflows</div>;
  },
  ForeignCardsSection: ({ cards }: any) => (
    <>{cards.map((card: any) => card.element)}</>
  ),
}));

function overview(workflows: boolean, observability = true) {
  return render(
    <TestApiProvider
      apis={[
        [
          configApiRef,
          mockApis.config({
            data: {
              openchoreo: {
                features: {
                  workflows: { enabled: workflows },
                  observability: { enabled: observability },
                },
              },
            },
          }),
        ],
      ]}
    >
      <ComponentOverview
        cards={[
          { element: <div key="about">About</div>, type: 'info' },
          { element: <div key="relations">Relations</div> },
          { element: <div key="third-party">Third-party card</div> },
        ]}
      />
    </TestApiProvider>,
  );
}

beforeEach(() => workflowMounted.mockClear());

test('disabled workflows never mount or request data', () => {
  overview(false);
  expect(screen.queryByText('Workflows')).not.toBeInTheDocument();
  expect(workflowMounted).not.toHaveBeenCalled();
  expect(screen.getByText('Deployment status')).toBeInTheDocument();
});

test('host About/Relations render once and third-party cards survive', () => {
  overview(false);
  expect(screen.getAllByText('About')).toHaveLength(1);
  expect(screen.getAllByText('Relations')).toHaveLength(1);
  expect(screen.getByText('Third-party card')).toBeInTheDocument();
  expect(
    screen.getByText(/not proof of deployment readiness/),
  ).toBeInTheDocument();
});

test('enabled workflows remain available', () => {
  overview(true);
  expect(screen.getByText('Workflows')).toBeInTheDocument();
  expect(workflowMounted).toHaveBeenCalledTimes(1);
});

test('disabled observability hides both its card and qualification', () => {
  overview(false, false);
  expect(screen.queryByText('Runtime health')).not.toBeInTheDocument();
  expect(
    screen.queryByText(/not proof of deployment readiness/),
  ).not.toBeInTheDocument();
});
