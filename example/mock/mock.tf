
data "ansible_inventory" "test" {
  hosts {
    group = "first"
    names = ["first-0", "first-1", "first-2"]
  }

  hosts {
    group = "second"
    names = ["second-0", "second-1"]

    vars {
      first = "first"
      second = <<EOF
      `cast:json` {
        "deep": 1,
        "deep_str": "`"
      }
      EOF
      ab = <<EOF
      `cast:json` {
        "b": true,
        "a": 1
      }
      EOF
    }
    vars {
      third = "`cast:string` '  third'"
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
  callback_whitelist = "profile_tasks"
}

resource "ansible_playbook" "test" {
  count = 2

  playbook = "${path.root}/ansible/playbook-1.yaml"
  inventory = "${data.ansible_inventory.test.rendered}"
  config = "${data.ansible_config.test_1.rendered}"
  phase {
    destroy = true
  }
  cleanup = false
  extra_json = <<EOF
  {
    "playbook_hash": "${sha256(file("${path.root}/ansible/playbook-1.yaml"))}",
    "a": 1,
    "b": true
  }
  EOF
}

