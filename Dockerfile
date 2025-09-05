# Start by building the application
FROM golang:1.24.4 AS build

WORKDIR /go/src

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /go/bin/xpense cmd/xpense/main.go

# Now copy it into our base image
FROM scratch AS runtime

WORKDIR /xpense

COPY --from=build /go/bin/xpense /xpense/xpense

EXPOSE 8090

ENTRYPOINT ["./xpense"]

CMD [ "serve", "--http", "0.0.0.0:8090" ]