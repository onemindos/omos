# Offline knowledge and home readiness

`/offline` provides an explicit private download of a single self-contained HTML file.
The backend requires an administrator user; anonymous users and machine credentials
cannot export it. Only the committed application directory and seven fixed packaged
runbooks are included. No live API results, access tokens, provider tasks, credentials
or conversation history are cached. There is no service worker or background cache.

The file uses no JavaScript, remote assets or automatic network requests. Content is
HTML-escaped, the browser supports local Find and printing, and a content security
policy disables scripts and forms. Account links work only when online and authorized.
Treat this as private internal documentation. Browser downloads cannot enforce filesystem
permissions or revocation; logout does not delete downloaded copies.

## Export without the portal

From a trusted local repository checkout, using Node:

```sh
node scripts/offline-pack.cjs --output /YOUR_PRIVATE_FOLDER/onemind-offline.html
```

Run from `apps/onemind-portal`. Output is explicitly selected, created with mode 0600
and never overwritten. The CLI refuses the application's public asset folder. The
snapshot is from that checkout, which may differ from the deployed image. Keep and
replace copies deliberately; don't synchronize them to public shared storage.

## What this does not provide

This PR does not install a home cluster or claim full offline application operation.
For that transition, follow the existing resilience runbook and choose actual compute,
storage and routing before rendering deployment manifests:

1. Preserve database backups, the decryption key, Git history and identity recovery data.
2. Choose a retained local storage class and restore PostgreSQL 17 into an isolated test
   instance before moving the active database.
3. Install local DNS, private HTTPS/Tailscale routing, identity and a registry/image mirror.
4. Re-seal credentials for the destination controller; cloud sealed secrets are not portable.
5. Repoint observations and catalog sources to the selected cluster and retain explicit
   context guards. Keep source timestamps and cloud dependency labels visible.
6. Rehearse disconnected read access and recovery on the actual destination hardware.

The output intentionally separates static knowledge from live functionality. Full
home deployment remains dependent on hardware and a verified restoration/cutover plan.
