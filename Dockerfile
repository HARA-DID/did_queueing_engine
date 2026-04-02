##############################################################################
# Stage 1 – build
##############################################################################
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Download dependencies first for layer caching.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /worker ./cmd/worker

##############################################################################
# Stage 2 – minimal runtime image
##############################################################################
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /worker /worker

EXPOSE 8080

ENTRYPOINT ["/worker"]
