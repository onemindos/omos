FROM python:3.11-slim
WORKDIR /app
RUN pip install --no-cache-dir clickhouse-connect
COPY nats-metrics-scrape.py .
CMD ["python", "-u", "nats-metrics-scrape.py"]
