terraform {
  required_version = "0.13.7"
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.36.0"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}


resource "aws_s3_bucket" "demo_bucket" {
  bucket = "${var.bucket_name}"
  versioning {
    enabled = true
  }
  tags = {
    enviornment   = "${var.tag_enviornment}"
    deployment    = "${var.tag_deployment}"
    region        = "${var.tag_region}"
  }
}

