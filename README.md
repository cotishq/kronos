# Kronos

Kronos is a small event-processing project I built to practice Kafka and gRPC together in Go.

The idea is simple: a client sends an event over gRPC, the gRPC server publishes it to Kafka, and separate consumers handle that event for analytics, audit logging, and dead-letter handling.

I built this in roughly two weeks, mostly as a hands-on learning project, with almost zero AI dependency. The hardest part was getting the gRPC side and the Kafka side calibrated properly so the flow felt clean instead of like two separate demos stitched together.

## Architecture

<img width="1238" height="714" alt="Screenshot From 2026-06-12 21-23-11" src="https://github.com/user-attachments/assets/26e01de5-bf89-4821-a9c2-b69f21beaefa" />


```text
grpc-client -> grpc-server -> Kafka topic: user-events
                              |-> analytics consumer
                              |-> audit consumer
                              |-> dead-letter topic -> dlq consumer

Prometheus scrapes:
- grpc-server metrics on :2112
- analytics-consumer metrics on :2113
```

## What It Does

- Accepts events through a gRPC service.
- Publishes events into Kafka.
- Processes events with separate Kafka consumer groups.
- Keeps a basic audit trail in `audit.log`.
- Tracks simple analytics counts per event type.
- Skips duplicate analytics events by event ID.
- Sends failed analytics events to a dead-letter topic after retries.
- Exposes Prometheus metrics for the gRPC server and analytics consumer.

## Project Structure

```text
cmd/
  grpc-server/          gRPC server that publishes incoming events to Kafka
  grpc-client/          small gRPC client for sending a test event
  producer/             direct Kafka producer with sample events
  consumer-analytics/   analytics consumer with retry and DLQ support
  consumer-audit/       audit consumer that writes to audit.log
  consumer-dlq/         consumer for the dead-letter topic

internal/event/         shared Event model
proto/                  protobuf definition and generated Go code
charts/                 Helm charts for Kubernetes deployment
  kronos/               main Helm chart for deployment
argocd/                 ArgoCD Application configuration
docker-compose.yml      Kafka, Prometheus, and Grafana for local Docker setup
prometheus.yml          Prometheus scrape config
```

## Tech Stack

- Go
- gRPC
- Protocol Buffers
- Kafka
- Prometheus
- Grafana
- Docker Compose
- Kubernetes
- Helm
- ArgoCD

## Running Locally

Start Kafka, Prometheus, and Grafana:

```bash
docker compose up -d
```

Run the gRPC server:

```bash
go run ./cmd/grpc-server
```

In another terminal, run the analytics consumer:

```bash
go run ./cmd/consumer-analytics
```

Run the audit consumer:

```bash
go run ./cmd/consumer-audit
```

Optionally run the DLQ consumer:

```bash
go run ./cmd/consumer-dlq
```

Send a test event through gRPC:

```bash
go run ./cmd/grpc-client
```

You can also publish the sample events directly to Kafka:

```bash
go run ./cmd/producer
```

## Kubernetes & ArgoCD Deployment

The project can be deployed to a Kubernetes cluster using the provided Helm chart or via ArgoCD for a GitOps workflow.

### 1. Deployment via Helm

To deploy the Kronos components using Helm:

```bash
# Deploy/Upgrade the Helm release
helm upgrade --install kronos ./charts/kronos
```

This will spin up:
- `grpc-server`
- `consumer-analytics`
- `consumer-audit`
- `consumer-dlq`
- `kafka` (local development cluster)

Deployment options (e.g. replicas, image tags, etc.) can be configured in [charts/kronos/values.yaml](charts/kronos/values.yaml).

### 2. GitOps Deployment via ArgoCD

An ArgoCD Application manifest is provided in [argocd/kronos-app.yml](argocd/kronos-app.yml).

To register the application in ArgoCD:

```bash
kubectl apply -f argocd/kronos-app.yml
```

This creates a `kronos` application tracking the `charts/kronos` path in this repository, with automated prune and self-healing enabled.

## Observability

The gRPC server exposes:

```text
http://localhost:2112/metrics
```

The analytics consumer exposes:

```text
http://localhost:2113/metrics
```

Prometheus is available at:

```text
http://localhost:9091
```

Grafana is available at:

```text
http://localhost:3001
```

## Notes

This is mainly a personal learning project. I built it to understand how Kafka consumers, gRPC APIs, retries, DLQs, and Prometheus metrics fit together in a small Go service.

Contributions are welcome if they improve the learning value of the project or make the implementation cleaner, but the project is not trying to become a production event platform.
