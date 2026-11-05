# ${{ values.title }}

Proposed by the Backstage "Add app" template. Pillar ${{ values.pillar }}. Namespace `${{ values.name }}`,
Argo CD Application `platform-${{ values.name }}` (manual sync).

## Runbook (fill in before the first sync)

- **What it does and who depends on it:**
- **Image and how it is built:** `${{ values.image }}`
- **Health:** {% if values.healthCheck == 'probe' %}readiness and liveness probe on `${{ values.healthPath }}`{% else %}none yet; add a probe or an `onemind.io/health-check` annotation{% endif %}
- **State and recovery:** {% if values.dataPlan == 'volsync' %}claim `data` on `do-block-storage-retain`, daily VolSync snapshot (`backup.yaml`); run an isolated restore before real use{% elif values.dataPlan == 'pg-dump' %}database; add a dump CronJob and verify one restore{% else %}no persistent state{% endif %}
- **Secrets:** see `${{ values.name }}.secret.example`; real values are SealedSecrets only.
- **Rollback:** revert the pull request that changed the image digest, then sync.
- **Retirement:** use the Backstage "Retire app" template; it lists everything attached and the manual teardown order.
