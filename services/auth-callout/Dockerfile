# OneMind NATS auth callout (FEATURES J2). Build context: services/ (same as services/policy).
#   docker build -f services/auth-callout/Dockerfile services
# Base images: the digests services/policy/Dockerfile pinned on 2026-10-29.
FROM golang:1.26.9-alpine@sha256:3082400e369fa24d5fc60bca20edab3f6d604e0c5a690ec66b295eff4dd87ade AS build
WORKDIR /src
COPY auth-callout/go.mod auth-callout/go.sum ./
RUN go mod download
COPY auth-callout/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth-callout ./cmd/auth-callout

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/auth-callout /auth-callout
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/auth-callout"]
