variable "aws_region" {
  description = "AWS region"
  default     = "us-east-1"
}

variable "instance_type" {
  description = "EC2 instance type"
  default     = "t3.nano"
}

variable "key_name" {
  description = "AWS SSH Key Pair name"
  default     = "your-key-name"
}

variable "ami_id" {
  description = "Amazon Machine Image (AMI) ID"
  default     = "ami-12345678" # Replace with latest Amazon Linux 2 AMI
}

variable "app_port" {
  description = "Port the app runs on"
  default     = 8080
}
