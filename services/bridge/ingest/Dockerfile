# syntax=docker/dockerfile:1
# omos-ingest image. Build context is services/ so the JSON Schemas in services/schemas are
# baked in at /schemas (see .github/workflows/build-omos-ingest.yml). Dockerfile.dockerignore
# limits the context to bridge/ingest and schemas.
# Pin both base images by digest after the first CI build (omos CLAUDE.md, deployment pattern 3).
FROM golang:1.27.3-alpine AS build
WORKDIR /src
COPY bridge/ingest/go.mod bridge/ingest/go.sum ./
RUN go mod download
COPY bridge/ingest/ ./
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/omos-ingest ./cmd/omos-ingest

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/omos-ingest /omos-ingest
COPY schemas/ /schemas/
USER 65532:65532
EXPOSE 9464
ENTRYPOINT ["/omos-ingest"]
