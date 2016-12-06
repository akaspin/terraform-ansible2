provider "aws" {
  region = "us-west-2"
}

resource "aws_vpc" "default" {
  cidr_block = "172.44.0.0/16"
  enable_dns_support = true
  enable_dns_hostnames = true
}

resource "aws_subnet" "default" {
  vpc_id = "${aws_vpc.default.id}"
  cidr_block = "172.44.1.0/24"
  availability_zone = "us-west-2a"
  map_public_ip_on_launch = true
}

resource "aws_internet_gateway" "default" {
  vpc_id = "${aws_vpc.default.id}"
}

resource "aws_route" "default" {
  route_table_id = "${aws_vpc.default.default_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${aws_internet_gateway.default.id}"
}

resource "aws_security_group" "default" {
  vpc_id = "${aws_vpc.default.id}"
  name = "example"
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_key_pair" "default" {
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}


resource "aws_instance" "first" {
  count = 2
  instance_type = "t2.nano"
  ami = "ami-d2c924b2"
  subnet_id = "${aws_subnet.default.id}"
  key_name = "${aws_key_pair.default.id}"
  vpc_security_group_ids = ["${aws_security_group.default.id}"]
  root_block_device {
    delete_on_termination = true
    volume_type = "gp2"
  }
  tags = {
    Name = "first-${count.index}"
  }
}

resource "aws_instance" "second" {
  count = 1
  instance_type = "t2.nano"
  ami = "ami-d2c924b2"
  subnet_id = "${aws_subnet.default.id}"
  key_name = "${aws_key_pair.default.id}"
  vpc_security_group_ids = ["${aws_security_group.default.id}"]
  root_block_device {
    delete_on_termination = true
    volume_type = "gp2"
  }
  tags = {
    Name = "second-${count.index}"
  }
}


