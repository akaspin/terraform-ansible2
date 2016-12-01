resource "ansible2_inventory" {
  hosts {
    group = "first"
    names = "${formatlist("%s", aws_instance.first.*.tags.Name)}"
  }
}

resource "ansible2_playbook" {
  
}
