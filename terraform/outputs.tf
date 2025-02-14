output "ec2_public_ip" {
  value = aws_instance.app_server.public_ip
}

output "elb_dns_name" {
  value = aws_elb.app_elb.dns_name
}
