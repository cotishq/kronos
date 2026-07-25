output "cluster_name" {
  value = module.eks.cluster_name
}

output "cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "cluster_certificate_authority_data" {
  value     = module.eks.cluster_certificate_authority_data
  sensitive = true
}

output "kubeconfig_command" {
  value = "aws eks update-kubeconfig --region ap-south-1 --name ${var.cluster_name}"
}