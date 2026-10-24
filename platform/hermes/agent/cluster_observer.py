#!/usr/bin/env python3
"""Collect bounded readiness evidence; the agent supplies analysis separately.

Only fixed GET list endpoints are used. Never emit raw API objects, pod specs,
environment values, condition/event messages, annotations, or credentials.
"""
import datetime
import json
from pathlib import Path
import ssl
import time
import urllib.error
import urllib.parse
import urllib.request

TOKEN_DIR = Path('/var/run/secrets/kubernetes.io/serviceaccount')
BASE_URL = 'https://kubernetes.default.svc'
MAX_BYTES = 16 * 1024 * 1024
MAX_PAGES = 20
TOTAL_SECONDS = 120
PATHS = {
    'nodes': '/api/v1/nodes',
    'namespaces': '/api/v1/namespaces',
    'pods': '/api/v1/pods',
    'services': '/api/v1/services',
    'endpoints': '/api/v1/endpoints',
    'events': '/api/v1/events',
    'persistentvolumeclaims': '/api/v1/persistentvolumeclaims',
    'persistentvolumes': '/api/v1/persistentvolumes',
    'deployments': '/apis/apps/v1/deployments',
    'statefulsets': '/apis/apps/v1/statefulsets',
    'daemonsets': '/apis/apps/v1/daemonsets',
    'replicasets': '/apis/apps/v1/replicasets',
    'jobs': '/apis/batch/v1/jobs',
    'cronjobs': '/apis/batch/v1/cronjobs',
    'ingresses': '/apis/networking.k8s.io/v1/ingresses',
    'networkpolicies': '/apis/networking.k8s.io/v1/networkpolicies',
    'endpointslices': '/apis/discovery.k8s.io/v1/endpointslices',
    'node_metrics': '/apis/metrics.k8s.io/v1beta1/nodes',
    'pod_metrics': '/apis/metrics.k8s.io/v1beta1/pods',
    'applications': '/apis/argoproj.io/v1alpha1/applications',
    'clusters': '/apis/postgresql.cnpg.io/v1/clusters',
    'replicationsources': '/apis/volsync.backube/v1alpha1/replicationsources',
}


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def conditions(status):
    # Messages may contain application input or credentials; never include them.
    return [{'type': c.get('type'), 'status': c.get('status'),
             'reason': c.get('reason')} for c in status.get('conditions', [])]


def summarize(kind, item):
    meta, spec, status = (item.get(k, {}) for k in ('metadata', 'spec', 'status'))
    row = {'namespace': meta.get('namespace'), 'name': meta.get('name'),
           'createdAt': meta.get('creationTimestamp')}
    if kind == 'nodes':
        row.update(conditions=conditions(status), allocatable=status.get('allocatable', {}),
                   capacity=status.get('capacity', {}), unschedulable=spec.get('unschedulable', False))
    elif kind == 'namespaces':
        row['phase'] = status.get('phase')
    elif kind == 'pods':
        containers = status.get('containerStatuses', []) + status.get('initContainerStatuses', [])
        row.update(phase=status.get('phase'), node=spec.get('nodeName'), conditions=conditions(status),
                   deleting=bool(meta.get('deletionTimestamp')),
                   resources=[{'name': c.get('name'),
                               'requests': {key: value for key, value in c.get('resources', {}).get('requests', {}).items()
                                            if key in ('cpu', 'memory', 'ephemeral-storage')},
                               'limits': {key: value for key, value in c.get('resources', {}).get('limits', {}).items()
                                          if key in ('cpu', 'memory', 'ephemeral-storage')}}
                              for c in spec.get('containers', [])],
                   containers=[{'name': c.get('name'), 'ready': c.get('ready'),
                                'restarts': c.get('restartCount'),
                                'state': next(iter(c.get('state', {})), None),
                                'reason': next(iter(c.get('state', {}).values()), {}).get('reason')}
                               for c in containers])
    elif kind in ('deployments', 'statefulsets', 'replicasets'):
        row.update(desired=spec.get('replicas', 1), ready=status.get('readyReplicas', 0),
                   updated=status.get('updatedReplicas', 0),
                   generation=meta.get('generation'), observedGeneration=status.get('observedGeneration'),
                   conditions=conditions(status))
    elif kind == 'daemonsets':
        row.update(desired=status.get('desiredNumberScheduled', 0), ready=status.get('numberReady', 0),
                   unavailable=status.get('numberUnavailable', 0), conditions=conditions(status))
    elif kind == 'jobs':
        row.update(active=status.get('active', 0), succeeded=status.get('succeeded', 0),
                   failed=status.get('failed', 0), conditions=conditions(status))
    elif kind == 'cronjobs':
        row.update(suspended=spec.get('suspend', False), active=len(status.get('active', [])),
                   lastScheduleTime=status.get('lastScheduleTime'),
                   lastSuccessfulTime=status.get('lastSuccessfulTime'))
    elif kind in ('persistentvolumeclaims', 'persistentvolumes'):
        row.update(phase=status.get('phase'), storageClass=spec.get('storageClassName'),
                   capacity=status.get('capacity', spec.get('capacity', {})))
    elif kind == 'services':
        row.update(type=spec.get('type'),
                   ports=[{'port': p.get('port'), 'protocol': p.get('protocol')}
                          for p in spec.get('ports', [])])
    elif kind == 'endpoints':
        row.update(readyAddresses=sum(len(s.get('addresses', [])) for s in item.get('subsets', [])),
                   notReadyAddresses=sum(len(s.get('notReadyAddresses', [])) for s in item.get('subsets', [])))
    elif kind == 'endpointslices':
        endpoints = item.get('endpoints', [])
        row.update(total=len(endpoints), ready=sum(e.get('conditions', {}).get('ready') is True for e in endpoints))
    elif kind == 'events':
        involved = item.get('involvedObject', {})
        row.update(type=item.get('type'), reason=item.get('reason'), count=item.get('count', 1),
                   lastTimestamp=item.get('lastTimestamp'),
                   involved={'kind': involved.get('kind'), 'name': involved.get('name')})
    elif kind == 'ingresses':
        row.update(ingressClass=spec.get('ingressClassName'),
                   loadBalancerReady=bool(status.get('loadBalancer', {}).get('ingress')))
    elif kind == 'networkpolicies':
        row.update(policyTypes=spec.get('policyTypes', []),
                   ingressRules=len(spec.get('ingress', [])), egressRules=len(spec.get('egress', [])))
    elif kind == 'node_metrics':
        row.update(usage=item.get('usage', {}), timestamp=item.get('timestamp'), window=item.get('window'))
    elif kind == 'pod_metrics':
        row.update(containers=[{'name': c.get('name'), 'usage': c.get('usage', {})}
                               for c in item.get('containers', [])], timestamp=item.get('timestamp'))
    elif kind == 'applications':
        row.update(sync=status.get('sync', {}).get('status'), health=status.get('health', {}).get('status'),
                   operationPhase=status.get('operationState', {}).get('phase'),
                   reconciledAt=status.get('reconciledAt'), conditions=conditions(status))
    elif kind == 'clusters':
        row.update(desired=spec.get('instances'), ready=status.get('readyInstances', 0),
                   phase=status.get('phase'), conditions=conditions(status))
    elif kind == 'replicationsources':
        row.update(lastSyncTime=status.get('lastSyncTime'), lastSyncDuration=status.get('lastSyncDuration'),
                   nextSyncTime=status.get('nextSyncTime'), conditions=conditions(status))
    else:
        raise ValueError('Unknown resource kind')
    return row


def collect(opener, deadline):
    sources = {}
    for kind, path in PATHS.items():
        rows, continuation = [], None
        try:
            for _ in range(MAX_PAGES):
                remaining = deadline - time.monotonic()
                if remaining <= 0:
                    raise TimeoutError()
                query = {'limit': '500'}
                if continuation:
                    query['continue'] = continuation
                request = urllib.request.Request(BASE_URL + path + '?' + urllib.parse.urlencode(query),
                                                 headers={'Accept': 'application/json',
                                                          'Authorization': 'Bearer ' + (TOKEN_DIR / 'token').read_text().strip()})
                with opener.open(request, timeout=min(10, remaining)) as response:
                    payload = response.read(MAX_BYTES + 1)
                if len(payload) > MAX_BYTES:
                    raise ValueError('Response size limit')
                body = json.loads(payload)
                rows.extend(summarize(kind, item) for item in body['items'])
                continuation = body.get('metadata', {}).get('continue')
                if not continuation:
                    break
            else:
                raise ValueError('Pagination limit')
            sources[kind] = {'complete': True, 'count': len(rows), 'items': rows}
        except Exception as exc:
            # Never print exception strings, HTTP bodies, request objects or auth.
            sources[kind] = {'complete': False, 'partialCount': len(rows), 'items': rows,
                             'error': type(exc).__name__}
            if isinstance(exc, urllib.error.HTTPError):
                sources[kind]['httpStatus'] = exc.code
    return sources


def main():
    started = datetime.datetime.now(datetime.timezone.utc).isoformat()
    context = ssl.create_default_context(cafile=str(TOKEN_DIR / 'ca.crt'))
    # Avoid ambient proxies and never redirect a bearer to a different host.
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}), NoRedirect(),
                                        urllib.request.HTTPSHandler(context=context))
    sources = collect(opener, time.monotonic() + TOTAL_SECONDS)
    print(json.dumps({'startedAt': started,
                      'completedAt': datetime.datetime.now(datetime.timezone.utc).isoformat(),
                      'complete': all(s['complete'] for s in sources.values()),
                      'scope': 'read-only readiness inventory; no logs or functional tests',
                      'sources': sources}, sort_keys=True))


if __name__ == '__main__':
    main()
