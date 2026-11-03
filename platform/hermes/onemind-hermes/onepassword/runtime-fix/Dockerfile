# Build context: this directory. Manual-only build-hermes-runtime.yml runs this gate.
FROM nousresearch/hermes-agent:v0.22.5@sha256:9c5561579f3ffdf333401797094ac3a743738d6c25d40f4331c7bf77b694238d
USER root
COPY apply_source_patch.py test_child_policy.py /opt/runtime-fix/
RUN /opt/hermes/.venv/bin/python /opt/runtime-fix/apply_source_patch.py --source-root /opt/hermes \
    && mkdir -p /opt/runtime-fix/scratch \
    && TMPDIR=/opt/runtime-fix/scratch \
       PYTHONPATH=/opt/hermes:/opt/hermes/.venv/lib/python3.14/site-packages \
       /opt/hermes/.venv/bin/python /opt/runtime-fix/test_child_policy.py --runtime-root /opt/hermes -v
# Preserve the upstream root entrypoint contract: s6 drops the gateway to UID 10000.
# No credential or mutable config is copied; manifests MUST use a real built digest.
