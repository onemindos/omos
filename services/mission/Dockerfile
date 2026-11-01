# OneMind task service (FEATURES D2, H4). Build context: services/mission.
#   docker build -f services/mission/Dockerfile services/mission
# Base images pinned by digest (the same digests as services/policy/Dockerfile, looked up
# 2026-10-29 from Docker Hub and gcr.io).
FROM golang:1.26.9-alpine@sha256:3082400e369fa24d5fc60bca20edab3f6d604e0c5a690ec66b295eff4dd87ade AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mission ./cmd/mission

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/mission /mission
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/mission"]
