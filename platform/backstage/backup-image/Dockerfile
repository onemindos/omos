FROM postgres:17-alpine@sha256:742f40ea20b9ff2ff31db5458d127452988a2164df9e17441e191f3b72252193
RUN apk add --no-cache python3 py3-boto3 py3-cryptography
COPY backup.py /backup.py
USER 70:70
ENTRYPOINT ["python3", "/backup.py"]
