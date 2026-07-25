module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = var.cluster_name
  cluster_version = var.eks_version

  vpc_id     = var.vpc_id
  subnet_ids = var.subnet_ids

  # public endpoint for kubectl access
  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  # cluster add-ons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
  }

  eks_managed_node_groups = {
    kronos = {
      instance_types = [var.node_instance_type]

      min_size     = var.node_min_count
      max_size     = var.node_max_count
      desired_size = var.node_desired_count

      # nodes live in private subnets
      subnet_ids = var.subnet_ids

      disk_size = 20  

      labels = {
        Project = "kronos"
      }

      tags = {
        Project     = "kronos"
        Environment = "dev"
        ManagedBy   = "terraform"
      }
    }
  }

  # allow kubectl from your machine
  enable_cluster_creator_admin_permissions = true

  tags = {
    Project     = "kronos"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

# IAM role for nodes to pull from ECR
resource "aws_iam_role_policy_attachment" "ecr_read" {
  role       = module.eks.eks_managed_node_groups["kronos"].iam_role_name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}