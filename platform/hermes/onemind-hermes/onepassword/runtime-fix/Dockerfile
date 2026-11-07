# Build context: this directory. Manual-only build-hermes-runtime.yml runs this gate.
FROM nousresearch/hermes-agent:v0.23.0@sha256:adb7890916de63245523492be8f68f9d2095fecbec4e2a6dab7960381f0fc002
USER root
COPY apply_source_patch.py test_child_policy.py /opt/runtime-fix/
RUN /opt/hermes/.venv/bin/python /opt/runtime-fix/apply_source_patch.py --source-root /opt/hermes \
    && mkdir -p /opt/runtime-fix/scratch \
    && TMPDIR=/opt/runtime-fix/scratch \
       PYTHONPATH=/opt/hermes:/opt/hermes/.venv/lib/python3.14/site-packages \
       /opt/hermes/.venv/bin/python /opt/runtime-fix/test_child_policy.py --runtime-root /opt/hermes -v
# Preserve the upstream root entrypoint contract: s6 drops the gateway to UID 10000.
# No credential or mutable config is copied; manifests MUST use a real built digest.
