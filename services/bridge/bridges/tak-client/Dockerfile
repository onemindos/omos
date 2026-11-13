# omos tak-client (TAK Server streaming client bridge). Base images pinned by digest
# (same digests as services/entity/Dockerfile on feat/entity-service, resolved 2026-10-29).
FROM golang:1.26-alpine@sha256:30b1e546a9c7fa3b747c2bdba3a0c7bdf3a02ebb4c5b8405d2f948d7e8aac4ec AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/tak-client .

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/tak-client /tak-client
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/tak-client"]
