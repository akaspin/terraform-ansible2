
data "ansible_inventory" "test" {
  hosts {
    group = "first"
    names = ["first-0", "first-1", "first-2"]
  }

  hosts {
    group = "second"
    names = ["second-0", "second-1"]

    var {
      key = "ansible_user1"
      value = "centos"
    }
  }

  var {
    group = "second"
    key = "public_ip"
    values = [
      "1.1.1.1",
      "2.2.2.2"
    ]
  }
}

data "ansible_config" "test_1" {
  roles_path = "./roles"
  task_includes_static = true
}

data "ansible_playbook" "test" {
  path = "${path.root}/ansible/playbook-1.yaml"
}

resource "ansible_playbook" "test" {
  count = 2

  inventory = "${data.ansible_inventory.test.rendered}"
  playbook = "${data.ansible_playbook.test.rendered}"
  config = "${data.ansible_config.test_1.rendered}"
  extra_json = <<EOF
  {
    "environment": "aws"
  }
  EOF
//  phase {
//    modify = true
//  }
//  phase {
//    modify = false
//  }
}

