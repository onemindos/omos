FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY meta_agent.py .
# Default: live heartbeat board. Override cmd for list/prompt/fanout.
ENTRYPOINT ["python", "meta_agent.py"]
CMD ["watch"]
