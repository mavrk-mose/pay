module "networking" {
  source = "./networking.tf"
}

module "compute" {
  source         = "./compute.tf"
  vpc_id        = module.networking.vpc_id
  public_subnet = module.networking.public_subnet
}

module "autoscaling" {
  source        = "./autoscaling.tf"
  launch_config = module.compute.launch_config
  public_subnet = module.networking.public_subnet
}

module "load_balancer" {
  source        = "./load_balancer.tf"
  public_subnet = module.networking.public_subnet
  ec2_instance  = module.compute.ec2_instance
}
