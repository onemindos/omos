# Cluster health observer

Argo application `platform-ops` owns this namespace and the five-minute observer.
It reads pod, volume and node state and prints structured findings. It never deletes
pods or repairs workloads. Pending pods can be legitimate during scheduling, and
failed Job pods must remain available for diagnosis until their controller/retention
policy removes them.

The service account is limited to get/list operations. The CronJob cannot overlap,
API calls have ten-second timeouts, and each Job has a two-minute deadline. Running
pods with unready containers are reported, including container waiting reasons.

Test the actual ConfigMap script with `python platform/ops/test_health_keeper.py`
(requires PyYAML). After reconciliation, create a Job from the CronJob, verify its
structured report, and confirm `kubectl auth can-i delete pods --as
system:serviceaccount:ops:health-keeper` returns no. The observer logs are evidence,
not an alert delivery channel; connect notifications separately through the platform's
approved event/monitoring path.
