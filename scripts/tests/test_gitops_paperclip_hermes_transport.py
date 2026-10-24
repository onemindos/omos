"""Security contracts for the existing-agent transport; no cluster credentials used."""
from pathlib import Path
import re
import shutil
import subprocess
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
HERMES = ROOT / 'platform/hermes/agent'
PAPERCLIP = ROOT / 'platform/paperclip'


def documents(path):
    return list(yaml.safe_load_all(path.read_text()))


class HermesTransport(unittest.TestCase):
    def test_tls_identities_and_private_service(self):
        docs = documents(HERMES / 'paperclip-api.yaml')
        ca = next(d for d in docs if d['kind'] == 'Certificate' and d['spec'].get('isCA'))
        self.assertEqual(ca['spec']['privateKey']['rotationPolicy'], 'Never')
        leaf = next(d for d in docs if d['kind'] == 'Certificate' and not d['spec'].get('isCA'))
        self.assertIn('hermes-paperclip-api.hermes.svc.cluster.local', leaf['spec']['dnsNames'])
        self.assertEqual(leaf['spec']['issuerRef']['name'], ca['metadata']['name'])
        svc = next(d for d in docs if d['kind'] == 'Service')
        self.assertEqual(svc['spec']['type'], 'ClusterIP')
        self.assertEqual(svc['spec']['ports'][0]['port'], 8643)
        self.assertNotIn('annotations', svc['metadata'])
        self.assertTrue(all(int(d['metadata']['annotations']['argocd.argoproj.io/sync-wave']) < 0
                            for d in docs if d['kind'] != 'Service'))

    def test_proxy_keeps_auth_and_streaming_and_denies_other_routes(self):
        config = (HERMES / 'paperclip-api-nginx.conf').read_text()
        self.assertIn('proxy_set_header Authorization $http_authorization;', config)
        self.assertIn('proxy_set_header Idempotency-Key $http_idempotency_key;', config)
        self.assertIn('proxy_set_header X-Hermes-Session-Key $http_x_hermes_session_key;', config)
        self.assertIn('proxy_buffering off;', config)
        self.assertIn('access_log off;', config)
        self.assertIn('location / { return 404; }', config)
        self.assertIn('ssl_certificate_cache max=2 inactive=10m valid=1m;', config)
        self.assertNotIn('auth_basic', config)
        self.assertNotIn('proxy_ssl_verify off', config)
        self.assertEqual(set(re.findall(r'proxy_pass ([^;]+);', config)), {'http://127.0.0.1:8642'})
        routes = re.findall(r'location [=~] (.+?) \{\n\s+limit_except ([A-Z ]+) \{ deny all; \}', config)
        # Evaluate exposed regexes against an adversarial route matrix.
        def allowed(path, method):
            return any((path == route if route.startswith('/') else re.fullmatch(route, path))
                       and method in methods.split() for route, methods in routes)
        for path, method in [('/health', 'GET'), ('/v1/capabilities', 'GET'),
                             ('/v1/runs', 'POST'), ('/v1/runs/abc-123', 'GET'),
                             ('/v1/runs/abc-123/events', 'GET'), ('/v1/runs/abc-123/stop', 'POST')]:
            self.assertTrue(allowed(path, method), (path, method))
        for path, method in [('/', 'GET'), ('/v1/skills', 'GET'), ('/v1/runs', 'DELETE'),
                             ('/v1/runs/abc/approval', 'POST'), ('/v1/runs/abc/steer', 'POST'),
                             ('/v1/runs/abc/stop', 'GET'), ('/health', 'POST'),
                             ('/v1/runs/abc/events/other', 'GET'), ('/v1/runs/../health', 'GET')]:
            self.assertFalse(allowed(path, method), (path, method))

    def test_certificate_path_exists_before_http_rewrite(self):
        # nginx -t does not perform a TLS handshake: a server `set $tls_dir`
        # passes syntax validation but is still unset during certificate lookup.
        # Require a lazy http-level map from the handshake-native SNI variable,
        # with a literal fixed result so untrusted SNI cannot select a file path.
        config = (HERMES / 'paperclip-api-nginx.conf').read_text()
        directives = re.sub(r'(?m)#.*$', '', config)
        http_start = directives.index('http {')
        server_start = directives.index('server {', http_start)
        before_server = directives[http_start:server_start]
        self.assertRegex(before_server,
                         r'map \$ssl_server_name \$tls_dir\s*\{\s*'
                         r'default /etc/hermes-paperclip-tls;\s*\}')
        self.assertNotRegex(directives, r'\bset\s+\$tls_dir\b')
        self.assertIn('ssl_certificate $tls_dir/tls.crt;', directives)
        self.assertIn('ssl_certificate_key $tls_dir/tls.key;', directives)
        self.assertNotRegex(directives, r'ssl_certificate(?:_key)?\s+[^;]*\$ssl_server_name')

    def test_sidecar_is_bounded_and_cannot_read_hermes_data(self):
        spec = documents(HERMES / 'deployment.yaml')[0]['spec']['template']['spec']
        sidecar = next(c for c in spec['containers'] if c['name'] == 'paperclip-api-tls')
        self.assertRegex(sidecar['image'], r'^nginx:1\.30\.6-alpine@sha256:[0-9a-f]{64}$')
        security = sidecar['securityContext']
        self.assertTrue(security['runAsNonRoot'])
        self.assertTrue(security['readOnlyRootFilesystem'])
        self.assertFalse(security['allowPrivilegeEscalation'])
        self.assertEqual(security['capabilities']['drop'], ['ALL'])
        self.assertEqual({v['name'] for v in sidecar['volumeMounts']},
                         {'paperclip-api-nginx', 'paperclip-api-tls', 'paperclip-api-tmp'})
        self.assertNotIn('envFrom', sidecar)
        self.assertEqual(sidecar['resources']['limits']['memory'], '64Mi')
        hermes = next(c for c in spec['containers'] if c['name'] == 'hermes')
        callback_mount = next(v for v in hermes['volumeMounts'] if v['name'] == 'paperclip-agent')
        self.assertEqual(callback_mount['mountPath'], '/run/paperclip')
        self.assertTrue(callback_mount['readOnly'])
        secret = next(v for v in spec['volumes'] if v['name'] == 'paperclip-agent')['secret']
        self.assertEqual(secret['items'], [{'key': 'credentials.json', 'path': 'credentials.json'}])
        self.assertEqual(secret['defaultMode'], 0o444)

    def test_paperclip_trust_and_exact_network_identity(self):
        deploy = next(d for d in documents(PAPERCLIP / 'deployment.yaml') if d['kind'] == 'Deployment')
        spec = deploy['spec']['template']['spec']
        env = spec['containers'][0]['env']
        self.assertIn({'name': 'NODE_EXTRA_CA_CERTS', 'value': '/run/hermes-gateway-ca/ca.crt'}, env)
        self.assertNotIn('NODE_TLS_REJECT_UNAUTHORIZED', [e['name'] for e in env])
        secret = next(v for v in spec['volumes'] if v['name'] == 'hermes-gateway-ca')['secret']
        self.assertEqual(secret['items'], [{'key': 'ca.crt', 'path': 'ca.crt'}])
        policy = next(d for d in documents(PAPERCLIP / 'network.yaml')
                      if d['metadata']['name'] == 'paperclip-network')['spec']
        peer = {'namespaceSelector': {'matchLabels': {'kubernetes.io/metadata.name': 'hermes'}},
                'podSelector': {'matchLabels': {'app.kubernetes.io/name': 'hermes',
                                               'onemind.io/agent': 'legacy'}}}
        self.assertIn({'to': [peer], 'ports': [{'protocol': 'TCP', 'port': 8643}]}, policy['egress'])
        self.assertIn({'from': [peer], 'ports': [{'protocol': 'TCP', 'port': 3100}]}, policy['ingress'])

    @unittest.skipUnless(shutil.which('kubectl'), 'kubectl kustomize not installed')
    def test_render_resolves_generated_configmap_and_namespaces(self):
        for directory, namespace in [(HERMES, 'hermes'), (PAPERCLIP, 'paperclip')]:
            result = subprocess.run(['kubectl', 'kustomize', str(directory)], check=True,
                                    capture_output=True, text=True)
            rendered = list(yaml.safe_load_all(result.stdout))
            self.assertTrue(all(d['metadata']['namespace'] == namespace for d in rendered
                                if d['kind'] not in ('Namespace', 'ClusterRole', 'ClusterRoleBinding')))
            if namespace == 'hermes':
                deploy = next(d for d in rendered if d['kind'] == 'Deployment')
                volume = next(v for v in deploy['spec']['template']['spec']['volumes']
                              if v['name'] == 'paperclip-api-nginx')
                self.assertTrue(any(d['kind'] == 'ConfigMap' and d['metadata']['name'] ==
                                    volume['configMap']['name'] for d in rendered))


if __name__ == '__main__':
    unittest.main()
