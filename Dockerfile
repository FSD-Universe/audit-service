FROM golang:1.25.4-alpine AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -tags "database http httpjwt grpc permission telemetry" -o /build/audit-service .

FROM alpine:latest AS runtime

RUN apk --no-cache add ca-certificates

WORKDIR /service

COPY --from=builder /build/audit-service .

ENTRYPOINT ["./audit-service"]