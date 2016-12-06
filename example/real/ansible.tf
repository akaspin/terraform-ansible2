
data "ansible_inventory" "test" {
  hosts {
    group = "first"
    names = ["${aws_instance.first.*.tags.Name}"]
  }
  var {
    group = "first"
    key = "ansible_host"
    values = ["${aws_instance.first.*.public_ip}"]
  }
  var {
    group = "first"
    key = "public_ip"
    values = ["${aws_instance.first.*.public_ip}"]
  }
}

data "ansible_playbook" "test" {
  path = "${path.root}/ansible/playbook-1.yaml"
}

data "ansible_config" "test" {
  remote_user = "centos"
  control_path = "/tmp/%%h-%%p-%%r"
  timeout = 300
}

resource "ansible_playbook" "test" {

  inventory = "${data.ansible_inventory.test.rendered}"
  playbook = "${data.ansible_playbook.test.rendered}"
  config = "${data.ansible_config.test.rendered}"
}
