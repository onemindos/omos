import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import type { EntityContentLayoutProps } from '@backstage/plugin-catalog-react/alpha';
import {
  DeploymentStatusCard,
  EntityWarningStrip,
  ForeignCardsSection,
  RuntimeHealthCard,
  WorkflowsOverviewCard,
} from '@openchoreo/backstage-plugin';
import { FeatureGate } from '@openchoreo/backstage-plugin-react';

// Upstream 1.3.1 embeds About/Relations and then renders the same host cards.
// Let Backstage own those cards once; retain other installed plugin cards too.
// This layout is selected only for OpenChoreo-managed Components.
export default function ComponentOverview({ cards }: EntityContentLayoutProps) {
  return (
    <Grid container spacing={3} alignItems="stretch">
      <EntityWarningStrip />
      <FeatureGate feature="workflows">
        <Grid item xs={12} md={4}>
          <WorkflowsOverviewCard />
        </Grid>
      </FeatureGate>
      <Grid item xs={12} md={4}>
        <DeploymentStatusCard />
      </Grid>
      <FeatureGate feature="observability">
        <Grid item xs={12} md={8}>
          <RuntimeHealthCard />
          <Typography variant="caption" component="p">
            Runtime health summarizes observed logs in the selected time window;
            no reported errors is not proof of deployment readiness or complete
            telemetry. Check Deploy and Metrics as well.
          </Typography>
        </Grid>
      </FeatureGate>
      <ForeignCardsSection cards={cards} />
    </Grid>
  );
}
