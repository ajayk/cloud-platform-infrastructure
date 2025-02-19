terraform {
  required_providers {
    auth0 = {
      source  = "alexkappa/auth0"
      version = "= 0.19.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.68.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.6.1"
    }
    null = {
      source = "hashicorp/null"
    }
  }
  required_version = ">= 0.14"
}
