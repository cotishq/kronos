module "vpc" {
  source       = "./modules/vpc"
  cluster_name = var.cluster_name
}

module "eks" {
  source             = "./modules/eks"
  cluster_name       = var.cluster_name
  eks_version        = var.eks_version
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnet_ids
  node_instance_type = var.node_instance_type
  node_desired_count = var.node_desired_count
  node_min_count     = var.node_min_count
  node_max_count     = var.node_max_count
}

module "ecr" {
  source       = "./modules/ecr"
  cluster_name = var.cluster_name
}