resource "aws_instance" "app_server" {
  ami           = var.ami_id
  instance_type = var.instance_type
  key_name      = var.key_name
  vpc_security_group_ids = [aws_security_group.app_sg.id]
  subnet_id     = aws_subnet.public_subnet.id

  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              sudo yum install -y docker
              sudo service docker start
              sudo usermod -aG docker ec2-user
              sudo amazon-linux-extras enable docker
              sudo systemctl enable docker
              sudo docker run -d -p ${var.app_port}:${var.app_port} your-docker-image
              EOF
  
  user_data = file("${path.module}/../scripts/setup-db.sh")

  tags = {
    Name = "payment-system-instance"
  }
}

resource "aws_elb" "app_elb" {
  name               = "payment-system-elb"
  availability_zones = ["us-east-1a"]

  listener {
    instance_port     = var.app_port
    instance_protocol = "HTTP"
    lb_port           = 80
    lb_protocol       = "HTTP"
  }

  instances = [aws_instance.app_server.id]

  tags = {
    Name = "payment-elb"
  }
}

resource "aws_autoscaling_group" "app_asg" {
  launch_configuration = aws_launch_configuration.app_lc.id
  min_size             = 1
  max_size             = 3
  vpc_zone_identifier  = [aws_subnet.public_subnet.id]

  tag {
    key                 = "Name"
    value               = "payment-system-instance"
    propagate_at_launch = true
  }
}

resource "aws_launch_configuration" "app_lc" {
  name          = "payment-lc"
  image_id      = var.ami_id
  instance_type = var.instance_type
  key_name      = var.key_name
  security_groups = [aws_security_group.app_sg.id]

  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              sudo yum install -y docker
              sudo service docker start
              sudo usermod -aG docker ec2-user
              sudo docker run -d -p ${var.app_port}:${var.app_port} your-docker-image
              EOF
}
