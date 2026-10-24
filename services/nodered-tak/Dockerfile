# Official Node-RED image; the multi-platform index is fixed, including its Node 24 runtime.
FROM nodered/node-red:5.0.9-24@sha256:e085d4a06591fd0aca2ea2a91364d27f1a0e31e885675f77935e10f95412a963
LABEL org.opencontainers.image.source="https://github.com/onemindos/omos" \
      org.opencontainers.image.description="Node-RED with locked TAK converters and authenticated safe-mode startup"

USER root
WORKDIR /opt/onemind-nodes
COPY package.json package-lock.json ./
RUN npm ci --omit=dev --omit=peer --omit=optional --ignore-scripts --no-audit --no-fund \
    && npm cache clean --force
COPY --chown=1000:1000 settings.cjs settings-factory.cjs bootstrap.cjs start.sh /opt/onemind/
COPY --chown=1000:1000 examples /opt/onemind/examples
RUN chmod 0555 /opt/onemind/start.sh
USER 1000:1000
WORKDIR /usr/src/node-red
ENTRYPOINT ["/opt/onemind/start.sh"]
CMD []
