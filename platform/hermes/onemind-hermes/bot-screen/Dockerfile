# Build context: this directory; publish through build-hermes-runtime.yml.
FROM nousresearch/hermes-agent:v0.22.5@sha256:9c5561579f3ffdf333401797094ac3a743738d6c25d40f4331c7bf77b694238d
USER root
# Only Xfce core: no session manager, locker, power manager or sudo.
RUN apt-get -o Acquire::Retries=3 update \
    && DEBIAN_FRONTEND=noninteractive apt-get -o Acquire::Retries=3 install -y --no-install-recommends \
       tigervnc-standalone-server xfce4-panel xfwm4 xfdesktop4 xfce4-settings \
       xfce4-terminal dbus-x11 x11-xserver-utils x11-utils x11-xkb-utils xauth fonts-dejavu-core \
    && apt-get clean
# PM is the upstream hash-locked authority. Do not install an unpinned driver via npm.
RUN cd /opt/hermes && /opt/hermes/.venv/bin/python -c 'from pm import ensure; ensure("cua-driver", explicit=True)'
COPY smoke_test.py /opt/bot-screen/smoke_test.py
# Preserve upstream s6 entrypoint: it drops the gateway to UID 10000.
# No desktop is started until requested; no credentials/config are baked in.
