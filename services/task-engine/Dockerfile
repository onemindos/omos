# omos-task-engine — Sovereign Task Orchestration Engine
# Multi-stage build: compile Go binary, run in scratch

FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /omos-task-engine .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /omos-task-engine /omos-task-engine

EXPOSE 4222
ENTRYPOINT ["/omos-task-engine"]
