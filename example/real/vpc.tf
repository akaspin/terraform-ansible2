module "vpc" {
  source = "../../../../../../../Craft/applift/borderline/terraform/modules/aws/vpc"

  site_name = "ansible"
  cidr = "172.34.0.0/16"
  region = "${var.aws_region}"

  subnet = {
    "172.34.0.0/24" = "a"
    "172.34.1.0/24" = "b"
    "172.34.2.0/24" = "c"
  }
}
