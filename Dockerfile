FROM node:22-alpine AS web-build

WORKDIR /app

COPY web/package.json web/pnpm-lock.yaml ./

RUN corepack enable pnpm && pnpm i --frozen-lockfile

COPY web .

RUN pnpm build

FROM golang:1.25.1 AS go-build

WORKDIR /go/src

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY web/embed.go ./web/embed.go

COPY --from=web-build /app/build ./web/build

RUN CGO_ENABLED=1 go build -o /go/bin/xpense cmd/xpense/main.go

FROM gcr.io/distroless/base-nossl AS runtime

WORKDIR /xpense

COPY --from=go-build /go/bin/xpense /xpense/xpense

EXPOSE 8080

CMD ["./xpense"]