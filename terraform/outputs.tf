output "ec2_public_ip" {
  value = aws_instance.app_server.public_ip
}

output "elb_dns_name" {
  value = aws_elb.app_elb.dns_name
}

output "database_endpoint" {
  value = "postgresql://payment_user:securepassword@${aws_instance.app_server.public_ip}:5432/payment_db"
  sensitive = true
}