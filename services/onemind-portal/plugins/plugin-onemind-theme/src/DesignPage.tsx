import {
  Button as MuiButton,
  Chip,
  FormControlLabel,
  Paper,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  makeStyles,
} from '@material-ui/core';
import { Button, Card, CardBody, CardHeader, Container } from '@backstage/ui';
import { contrastPairs, contrastRatio } from './contrast';
import { StatusChip } from './StatusChip';
import { useBackgroundImage } from './preferences';
import {
  categoryPalette,
  cssVariables,
  CSS_PREFIX,
  StatusKey,
  ThemeMode,
  tokens,
} from './tokens';

const useStyles = makeStyles({
  grid: {
    display: 'grid',
    gap: 16,
    gridTemplateColumns: 'repeat(auto-fit,minmax(min(100%,420px),1fr))',
  },
  swatches: {
    display: 'grid',
    gap: 8,
    gridTemplateColumns: 'repeat(auto-fill,minmax(120px,1fr))',
  },
  swatch: {
    border: '1px solid var(--onemind-border)',
    borderRadius: 8,
    overflow: 'hidden',
    fontSize: 12,
  },
  chip: { height: 44 },
  caption: { padding: '4px 8px', color: 'var(--onemind-text-2)' },
  row: { display: 'flex', flexWrap: 'wrap', gap: 8, alignItems: 'center' },
});

const STATUSES: StatusKey[] = ['ok', 'warn', 'critical', 'unknown', 'info'];

function Swatches({ mode }: { mode: ThemeMode }) {
  const classes = useStyles();
  const vars = cssVariables(mode);
  return (
    <div className={classes.swatches}>
      {Object.entries(vars).map(([name, value]) => (
        <div key={name} className={classes.swatch}>
          <div className={classes.chip} style={{ background: value }} />
          <div className={classes.caption}>
            <code>{name.replace(CSS_PREFIX, '')}</code>
            <br />
            {value}
          </div>
        </div>
      ))}
    </div>
  );
}

function Contrast({ mode }: { mode: ThemeMode }) {
  const pairs = contrastPairs(mode);
  return (
    <Table size="small">
      <TableHead>
        <TableRow>
          <TableCell>Pair</TableCell>
          <TableCell>Ratio</TableCell>
          <TableCell>Needs</TableCell>
          <TableCell>Result</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {pairs.map(p => {
          const ratio = contrastRatio(p.fg, p.bg);
          return (
            <TableRow key={p.name}>
              <TableCell>{p.name}</TableCell>
              <TableCell>{ratio.toFixed(2)}:1</TableCell>
              <TableCell>{p.min}:1</TableCell>
              <TableCell>
                <StatusChip
                  status={ratio >= p.min ? 'ok' : 'critical'}
                  label={ratio >= p.min ? 'Pass' : 'Fail'}
                />
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}

/** `/design`: every token and the main components, for palette reviews. */
export function DesignPage() {
  const classes = useStyles();
  const [background, setBackground] = useBackgroundImage();
  return (
    <Container my="6">
      <Typography variant="h4" component="h1" gutterBottom>
        OneMind design system
      </Typography>
      <Typography color="textSecondary" paragraph>
        Tokens live in <code>plugin-onemind-theme/src/tokens.ts</code> and reach
        MUI, Backstage UI and plain CSS (<code>--onemind-*</code>). The theme
        you are viewing is applied to the live components below; the swatches
        show both modes.
      </Typography>
      <FormControlLabel
        control={
          <Switch
            checked={background}
            onChange={(_e, v) => setBackground(v)}
            color="primary"
          />
        }
        label="Show the Legacy background image behind page headers (off by default)"
      />

      <Typography
        variant="h5"
        component="h2"
        gutterBottom
        style={{ marginTop: 24 }}
      >
        Components
      </Typography>
      <div className={classes.grid}>
        <Paper variant="outlined" style={{ padding: 16 }}>
          <Typography variant="h6" component="h3" gutterBottom>
            Buttons and chips (MUI)
          </Typography>
          <div className={classes.row}>
            <MuiButton variant="contained" color="primary">
              Primary
            </MuiButton>
            <MuiButton variant="outlined" color="primary">
              Outlined
            </MuiButton>
            <MuiButton variant="text" color="primary">
              Text
            </MuiButton>
            <MuiButton variant="contained" color="secondary">
              Accent
            </MuiButton>
            <MuiButton variant="contained" disabled>
              Disabled
            </MuiButton>
          </div>
          <div className={classes.row} style={{ marginTop: 12 }}>
            <Chip label="Neutral chip" />
            <Chip label="Outlined chip" variant="outlined" />
            <a href="#links">A link</a>
          </div>
        </Paper>
        <Card>
          <CardHeader>Buttons and cards (Backstage UI)</CardHeader>
          <CardBody>
            <div className={classes.row}>
              <Button variant="primary">Primary</Button>
              <Button variant="secondary">Secondary</Button>
              <Button variant="tertiary">Tertiary</Button>
            </div>
          </CardBody>
        </Card>
        <Paper variant="outlined" style={{ padding: 16 }}>
          <Typography variant="h6" component="h3" gutterBottom>
            Status (health only)
          </Typography>
          <div className={classes.row}>
            {STATUSES.map(k => (
              <StatusChip key={k} status={k} />
            ))}
          </div>
          <Typography
            variant="body2"
            color="textSecondary"
            style={{ marginTop: 8 }}
          >
            Critical is a rose red with an alert icon; the brand red is never
            used for health.
          </Typography>
        </Paper>
        <Paper variant="outlined" style={{ padding: 16 }}>
          <Typography variant="h6" component="h3" gutterBottom>
            Dense table
          </Typography>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Service</TableCell>
                <TableCell>Namespace</TableCell>
                <TableCell>Health</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {[
                ['argocd', 'argocd', 'ok'],
                ['backstage', 'backstage', 'ok'],
                ['loki', 'observability', 'warn'],
                ['harbor', 'registry', 'critical'],
                ['gatus', 'monitoring', 'unknown'],
              ].map(([n, ns, st]) => (
                <TableRow key={n}>
                  <TableCell>{n}</TableCell>
                  <TableCell>{ns}</TableCell>
                  <TableCell>
                    <StatusChip status={st as StatusKey} />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Paper>
        <Paper variant="outlined" style={{ padding: 16 }}>
          <Typography variant="h6" component="h3" gutterBottom>
            Typography
          </Typography>
          <Typography variant="h5">Heading five</Typography>
          <Typography variant="h6">Heading six</Typography>
          <Typography>Body text, primary.</Typography>
          <Typography color="textSecondary">Secondary text.</Typography>
          <Typography variant="caption">Caption text.</Typography>
        </Paper>
        <Paper variant="outlined" style={{ padding: 16 }}>
          <Typography variant="h6" component="h3" gutterBottom>
            Category hues (decorative)
          </Typography>
          <div className={classes.row}>
            {Object.entries(categoryPalette).map(([name, c]) => (
              <span
                key={name}
                style={{ borderLeft: `4px solid ${c.base}`, paddingLeft: 8 }}
              >
                {name}
              </span>
            ))}
          </div>
        </Paper>
      </div>

      {(['dark', 'light'] as const).map(mode => (
        <section key={mode}>
          <Typography
            variant="h5"
            component="h2"
            gutterBottom
            style={{ marginTop: 32 }}
          >
            Tokens: {mode}
          </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            Surface 0 {tokens[mode].surface[0]}, brand {tokens[mode].brand.base}
            .
          </Typography>
          <Swatches mode={mode} />
          <Typography
            variant="h6"
            component="h3"
            gutterBottom
            style={{ marginTop: 16 }}
          >
            WCAG contrast: {mode}
          </Typography>
          <Contrast mode={mode} />
        </section>
      ))}
    </Container>
  );
}
