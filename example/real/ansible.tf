
data "ansible_inventory" "test" {
  hosts {
    group = "first"
    names = "${join(",", module.instance_first.hostname)}"
  }
  var {
    group = "first"
    key = "ansible_host"
    values = "${join(",", module.instance_first.public_ip)}"
  }
}


//resource "null_resource" "test" {
//  triggers = {
//    aaa = "${data.ansible_inventory.test.rendered}"
//  }
//}


data "ansible_playbook" "test" {
  path = "${path.root}/ansible/playbook-1.yaml"
}

data "ansible_config" "test" {
  remote_user = "centos"
  control_path = "/tmp/%%h-%%p-%%r"
}

resource "ansible_playbook" "test" {
  inventory = "${data.ansible_inventory.test.rendered}"
  playbook = "${data.ansible_playbook.test.rendered}"
  config = "${data.ansible_config.test.rendered}"
}
