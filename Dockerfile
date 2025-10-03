FROM node:22-alpine AS web-build

WORKDIR /app

COPY web/package.json web/pnpm-lock.yaml ./

RUN corepack enable pnpm && pnpm i --frozen-lockfile

COPY web .

RUN pnpm build

FROM golang:1.25.1 AS go-dependencies

WORKDIR /go/src

COPY go.mod go.sum ./

RUN go mod download

FROM go-dependencies AS go-build-healthcheck

COPY cmd/healthcheck ./cmd/healthcheck
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /go/bin/healthcheck cmd/healthcheck/main.go

FROM go-dependencies AS go-build-main

COPY cmd/xpense ./cmd/xpense
COPY internal ./internal
COPY web/embed.go ./web/embed.go

COPY --from=web-build /app/build ./web/build

RUN CGO_ENABLED=1 go build -o /go/bin/xpense cmd/xpense/main.go

FROM gcr.io/distroless/base-nossl AS runtime

WORKDIR /xpense

COPY --from=go-build-healthcheck /go/bin/healthcheck /xpense/healthcheck
COPY --from=go-build-main /go/bin/xpense /xpense/xpense

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD ["/xpense/healthcheck"]

EXPOSE 8080

CMD ["./xpense"]