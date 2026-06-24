FROM golang:1.26.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/grpc-server ./cmd/grpc-server

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/bin/grpc-server .

EXPOSE 50051 2112
USER nonroot:nonroot

ENTRYPOINT [ "/app/grpc-server" ]