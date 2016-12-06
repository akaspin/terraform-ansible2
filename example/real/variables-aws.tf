variable "aws_region" {
  type = "string"
}

module "instance_first" {
  source = "../../../../../../../Craft/applift/borderline/terraform/modules/aws/instance"

  number = "3"

  site_name = "ansible"
  name = "proxy"
  cidr = "172.34.0.0/16"

  aws_vpc_id = "${module.vpc.vpc_id}"
  aws_region = "${var.aws_region}"
  aws_subnet_ids = "${module.vpc.subnet_ids}"
  aws_instance_type = "t2.small"
}
