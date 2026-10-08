import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch, MagicMock
import backup


class KubernetesRecoveryTests(unittest.TestCase):
    def test_paginated_export_preserves_sealing_keys_and_strips_server_metadata(self):
        secret = {'metadata': {'name': 'key', 'namespace': 'kube-system', 'uid': 'old-id',
                               'annotations': {'last-applied': 'secret'},
                               'labels': {'sealedsecrets.bitnami.com/sealed-secrets-key': 'active'}},
                  'type': 'Opaque', 'data': {'tls.key': 'synthetic'}}
        pages = [{'items': [secret], 'metadata': {'continue': 'next-page'}},
                 {'items': [], 'metadata': {}},
                 {'items': [{'metadata': {'name': 'config', 'namespace': 'nats'},
                             'data': {'config': 'synthetic'}}], 'metadata': {}}]
        responses = []
        for page in pages:
            response = MagicMock()
            response.__enter__.return_value.read.return_value = json.dumps(page)
            responses.append(response)
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            (root / 'token').write_text('synthetic-token')
            with patch.dict(os.environ, {'KUBE_TOKEN_DIR': tmp, 'KUBERNETES_SERVICE_HOST': 'synthetic'}), \
                    patch.object(backup.ssl, 'create_default_context'), \
                    patch.object(backup.urllib.request, 'urlopen', side_effect=responses) as opened:
                metadata = backup.kubernetes_backup(root)
            self.assertEqual(opened.call_count, 3)
            self.assertIn('continue=next-page', opened.call_args_list[1].args[0].full_url)
            record = json.loads((root / 'recovery.json').read_text())
            self.assertNotIn('uid', record['items'][0]['metadata'])
            self.assertNotIn('annotations', record['items'][0]['metadata'])
            self.assertEqual(metadata['sealingKeys'], 1)
            self.assertTrue(backup.kubernetes_verify(root, metadata)['downloadDecryptEqualityVerified'])
            (root / 'recovery.json').write_text('{}')
            with self.assertRaisesRegex(RuntimeError, 'equality mismatch'):
                backup.kubernetes_verify(root, metadata)


if __name__ == '__main__':
    unittest.main()
