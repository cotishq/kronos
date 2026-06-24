FROM golang:1.26.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/producer ./cmd/producer

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/bin/producer .

USER nonroot:nonroot

ENTRYPOINT [ "/app/producer" ]