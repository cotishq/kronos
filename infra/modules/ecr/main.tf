locals {
  services = ["grpc-server", "consumer-analytics", "consumer-audit", "consumer-dlq", "producer"]
}

resource "aws_ecr_repository" "kronos" {
  for_each = toset(local.services)

  name                 = "kronos/${each.key}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Project     = "kronos"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

resource "aws_ecr_lifecycle_policy" "kronos" {
  for_each   = aws_ecr_repository.kronos
  repository = each.value.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}