import { useEntity } from '@backstage/plugin-catalog-react';
import { InfoCard, Link } from '@backstage/core-components';
import { Box, Typography } from '@material-ui/core';
import { InfoTooltip } from '@internal/plugin-onemind-common';
import { buildFacts, Fact, NOT_RECORDED } from './facts';

function renderFact(fact: Fact) {
  if (fact.value === undefined) return <em style={{ opacity: 0.7 }}>{NOT_RECORDED}</em>;
  if (fact.link) return <Link to={fact.link}>{fact.value}</Link>;
  return fact.value;
}

/** Owner, hosting, plan, renewal, export, backup, recovery and access. Missing facts say "Not recorded". */
export function FactsCard() {
  const { entity } = useEntity();
  const facts = buildFacts(entity);
  return (
    <InfoCard
      title="App facts"
      action={
        <InfoTooltip label="app facts">
          Text recorded in the catalog (onemind.io annotations). Nothing here is checked live and nothing is guessed:
          a blank field means nobody has recorded it. Passwords and keys are never stored in the catalog; the access line
          says who can sign in, not how.
        </InfoTooltip>
      }
    >
      <Box component="dl" m={0} display="grid" style={{ gridTemplateColumns: 'max-content 1fr', columnGap: 16, rowGap: 6 }}>
        {facts.map(fact => (
          <Box key={fact.label} display="contents">
            <Typography component="dt" variant="body2" color="textSecondary">{fact.label}</Typography>
            <Typography component="dd" variant="body2" style={{ margin: 0 }} data-testid={`fact-${fact.label}`}>
              {renderFact(fact)}
            </Typography>
          </Box>
        ))}
      </Box>
    </InfoCard>
  );
}
