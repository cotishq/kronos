# Kronos

A distributed event processing pipeline built in Go - gRPC ingestion, Kafka fan-out,
multiple consumers, deployed on Kubernetes with full observability as code.

## Architecture


<img width="1103" height="714" alt="image" src="https://github.com/user-attachments/assets/db6ba787-7bf3-4774-beaa-215f7d732c1d" />



```text
grpc-client → grpc-server → Kafka topic: user-events
                                  |→ consumer-analytics
                                  |→ consumer-audit
                                  |→ dead-letter topic → consumer-dlq

Kubernetes (kind / EKS):
  - Helm chart → ArgoCD GitOps → auto-synced on push
  - kube-prometheus-stack → Prometheus + Grafana + Alertmanager
  - ServiceMonitor → auto-discovers services
  - PrometheusRules → burn-rate alerting on DLQ spike
```

## Stack

| Layer | Tech |
|-------|------|
| API | gRPC + Protocol Buffers |
| Messaging | Kafka |
| Language | Go |
| Containers | Docker |
| Orchestration | Kubernetes (kind locally, EKS on AWS) |
| Package Management | Helm |
| GitOps | ArgoCD |
| Observability | Prometheus, Grafana, Alertmanager |
| IaC | Terraform (VPC + EKS + ECR modules) |

## Project Structure

```text
cmd/
  grpc-server/         gRPC server, publishes events to Kafka
  grpc-client/         test client
  producer/            direct Kafka producer
  consumer-analytics/  analytics consumer with retry + DLQ
  consumer-audit/      audit consumer → audit.log
  consumer-dlq/        dead-letter consumer

charts/kronos/         Helm chart
  templates/
    grafana-dashboard.yaml    Grafana dashboard provisioned as ConfigMap
    prometheusrule.yaml       PrometheusRules for burn-rate alerting
    servicemonitor.yaml       ServiceMonitor for auto-discovery

argocd/                ArgoCD Application manifest
infra/                 Terraform modules
  modules/vpc/         VPC + subnets + NAT
  modules/eks/         EKS cluster + managed node group
  modules/ecr/         ECR repos per service + lifecycle policies
```

## Observability

kube-prometheus-stack runs inside the cluster — Prometheus, Grafana, and Alertmanager
all deployed via Helm, auto-discovering services via ServiceMonitor.

Dashboards and alerting rules are committed as code — no manual Grafana clicking.

**Metrics exposed:**
- `kronos_events_published_total` — events into Kafka
- `kronos_events_processed_total` — events consumed by analytics
- `kronos_events_dlq_total` — failed events hitting DLQ

**Alerts:**
- `KronosDLQSpike` — fires when DLQ rate > 0 for 2 minutes

```bash
# access Grafana
kubectl port-forward -n monitoring svc/monitoring-grafana 3001:80
# http://localhost:3001 → admin / <grafana.adminPassword from values.yaml>
```

## Infrastructure (Terraform)

Terraform modules provision the full AWS stack:

```bash
cd infra
terraform init
terraform plan
terraform apply
# then:
aws eks update-kubeconfig --region ap-south-1 --name kronos-eks
```

Remote state stored in S3 with DynamoDB locking.

## Local Setup

```bash
# start Kafka + Zookeeper
docker compose up -d

# run services
go run ./cmd/grpc-server
go run ./cmd/consumer-analytics
go run ./cmd/consumer-audit
go run ./cmd/consumer-dlq

# send a test event
go run ./cmd/grpc-client -addr=localhost:50051
```

## Kubernetes Deployment

```bash
# via Helm
helm upgrade --install kronos ./charts/kronos

# via ArgoCD (GitOps)
kubectl apply -f argocd/kronos-app.yml
# ArgoCD auto-syncs on every push to main
```
