terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.68.0"
    }
    github = {
      source = "integrations/github"
    }
  }
  required_version = ">= 0.14"
}
