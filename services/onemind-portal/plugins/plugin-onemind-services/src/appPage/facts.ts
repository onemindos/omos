import { hostingKind, HOSTING_KIND_LABEL } from '@internal/plugin-onemind-common';
import { RUNBOOK_SITE } from './links';

type EntityLike = {
  metadata: { annotations?: Record<string, string> };
  spec?: Record<string, unknown>;
};

export const NOT_RECORDED = 'Not recorded';

export type Fact = { label: string; value?: string; link?: string };

const text = (a: Record<string, string>, key: string) => {
  const v = a[key]?.trim();
  return v ? v : undefined;
};

/**
 * The facts block on the app page. Every field is text from the catalog (onemind.io/* annotations, see
 * docs/catalog/app-contract.md); a field nobody recorded has no value and the card prints "Not recorded".
 * Nothing is inferred, and credentials are never stored here.
 */
export function buildFacts(entity: EntityLike): Fact[] {
  const a = entity.metadata.annotations ?? {};
  const owner = typeof entity.spec?.owner === 'string' ? entity.spec.owner.replace(/^group:(default\/)?/, '') : undefined;
  const kind = hostingKind(entity);
  const runbook = text(a, 'onemind.io/recovery-runbook');
  const isUrl = !!runbook && /^https?:\/\//.test(runbook);
  const recovery: Fact = { label: 'Recovery' };
  if (runbook) {
    recovery.value = isUrl ? runbook : `Runbook ${runbook}`;
    recovery.link = isUrl ? runbook : `${RUNBOOK_SITE}${runbook}`;
  }
  return [
    { label: 'Owner', value: owner },
    {
      label: 'Hosting',
      value: kind === 'unknown' ? 'Unknown: hosting has not been established' : HOSTING_KIND_LABEL[kind],
    },
    { label: 'Plan', value: text(a, 'onemind.io/plan') },
    { label: 'Renewal', value: text(a, 'onemind.io/renewal') },
    { label: 'Data export', value: text(a, 'onemind.io/data-export') },
    { label: 'Backup location', value: text(a, 'onemind.io/backup-location') },
    recovery,
    { label: 'Who has access', value: text(a, 'onemind.io/access') },
  ];
}
