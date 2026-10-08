ARG POSTGRES_IMAGE=postgres:15-alpine@sha256:f7d23353e1b15400d22ebe31189f4d314b87a4c129cc400c8c2d8d4ca127bf81
FROM natsio/nats-box:0.20.0@sha256:d2f04223f40e016f77922a8115d34b0cd5fad2ec10b19b39316ae92d07f3d295 AS cli
FROM nats:2.14.2-alpine@sha256:952d157e28d5394a211229bd57a7b37ff9f184e58e2c8486a08fa909fd254e32 AS server
FROM ${POSTGRES_IMAGE}
RUN apk add --no-cache python3 py3-boto3 py3-cryptography
COPY --from=cli /usr/local/bin/nats /usr/local/bin/nats
COPY --from=server /nats-server /usr/local/bin/nats-server
COPY backup.py /backup.py
ENV HOME=/tmp
USER 70:70
ENTRYPOINT ["python3", "/backup.py"]
