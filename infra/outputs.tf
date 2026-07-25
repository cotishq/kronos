output "eks_cluster_name" {
  value = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "kubeconfig_command" {
  value = module.eks.kubeconfig_command
}

output "ecr_repository_urls" {
  value = module.ecr.repository_urls
}

output "vpc_id" {
  value = module.vpc.vpc_id
}