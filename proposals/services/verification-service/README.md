# verification-service proposal

Owner: OneMind administrators. This directory is not included in GitOps automatically.

Before deployment, review the pinned image, non-root UID, /health readiness/liveness behavior, resource sizing, Namespace, explicit network ingress/egress and private mesh routing. The default policy denies all network traffic. Add only required peers and destinations, including DNS if needed.

Add approved Secret references to envFrom/env; create and seal those Secrets separately. Never put values in this proposal. No credentials are needed by default.

Add the service to the catalog, monitoring and a runbook. Confirm monitoring endpoints; the template does not claim the image exposes Prometheus metrics. Review its HTTP health path and include it in the supported monitoring system before deployment.

A reviewer must move approved files into platform/, add Namespace and kustomization, and create the corresponding Argo Application. Validate with a server-side dry run and review before merge. Roll back through Git; retain persistent storage where applicable. Repeated runs update the same proposal branch and PR.
