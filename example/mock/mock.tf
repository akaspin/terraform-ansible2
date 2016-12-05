data "ansible2_inventory" "test" {
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

data "ansible2_config" "test_1" {
  roles_path = "./roles"
  task_includes_static = true
}

resource "ansible2_play" "test" {
  count = 2
  
  inventory = "${data.ansible2_inventory.test.rendered}"
  playbook = "${path.root}/ansible/playbook-1.yaml"
//  config = "${data.ansible2_config.test_1.rendered}"
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

