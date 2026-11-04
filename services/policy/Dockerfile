# OneMind policy and audit service (FEATURES I4). Build context: services/ (so the image
# carries the action, task and envelope schemas from services/schemas, the single source).
#   docker build -f services/policy/Dockerfile services
# Base images are pinned by digest (looked up 2026-10-29 from Docker Hub and gcr.io).
FROM golang:1.26.9-alpine@sha256:3082400e369fa24d5fc60bca20edab3f6d604e0c5a690ec66b295eff4dd87ade AS build
WORKDIR /src
COPY policy/go.mod policy/go.sum ./
RUN go mod download
COPY policy/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/policy ./cmd/policy

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/policy /policy
COPY schemas/actions /schemas/actions
COPY schemas/tasks /schemas/tasks
COPY schemas/command /schemas/command
COPY schemas/envelope/envelope.schema.json /schemas/envelope/envelope.schema.json
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/policy"]
